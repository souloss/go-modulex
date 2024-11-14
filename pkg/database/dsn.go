package database

import (
	"strings"
)

type DataSource struct {
	Scheme    string `mapstructure:"type" `
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Transport string `mapstructure:"transport"`
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Database  string `mapstructure:"db"`
	SSLMode   string `mapstructure:"ssl_mode"`

	AdditionalParams map[string]string `mapstructure:"additional_params"`

	Raw string `mapstructure:"dsn"`
}

func (ds *DataSource) GenerateDSN() string {
	var dsn strings.Builder
	dsn.WriteString(ds.Scheme) // 数据库类型
	dsn.WriteString("://")

	if ds.User != "" {
		dsn.WriteString(ds.User)
		if ds.Password != "" {
			dsn.WriteString(":" + ds.Password)
		}
		dsn.WriteString("@")
	}

	if ds.Host != "" {
		dsn.WriteString(ds.Host)
	}

	if ds.Port != "" {
		dsn.WriteString(":" + ds.Port)
	}

	if ds.Database != "" {
		dsn.WriteString("/") // 数据库名
		dsn.WriteString(ds.Database)
	}

	if len(ds.AdditionalParams) > 0 || ds.SSLMode != "" {
		dsn.WriteString("?") // 查询参数开始
		first := true
		for key, value := range ds.AdditionalParams {
			if !first {
				dsn.WriteString("&")
			}
			dsn.WriteString(key)
			dsn.WriteString("=")
			dsn.WriteString(value)
			first = false
		}
		if ds.SSLMode != "" {
			if !first {
				dsn.WriteString("&")
			}
			dsn.WriteString("sslmode")
			dsn.WriteString("=")
			dsn.WriteString(ds.SSLMode)
		}
	}

	return dsn.String()
}

// Parse receives a raw dsn as argument, parses it and returns it in the DSN structure.
func Parse(raw string) *DataSource {
	d := DataSource{
		Raw:              raw,
		AdditionalParams: map[string]string{},
	}
	dsn := []rune(d.Raw)

	// Parsing the scheme
	for pos, symbol := range dsn {
		// Found end of the scheme name
		if symbol == '/' && pos > 2 && string(dsn[pos-2:pos+1]) == "://" {
			d.Scheme = string(dsn[0 : pos-2])
			dsn = dsn[pos+1:]
			break
		}
	}

	// Parsing the credentials
	for dsnPos, dsnSymbol := range dsn {
		// Found end of the credentials
		if dsnSymbol == '@' && !isEscaped(dsnPos, dsn) {
			credentials := dsn[0:dsnPos]

			// Separating username and password
			hasSeparator := false
			for credPos, credChar := range credentials {
				if credChar == ':' && !isEscaped(credPos, credentials) {
					hasSeparator = true
					d.User = string(unEscape([]rune{':', '@'}, credentials[0:credPos]))
					d.Password = string(unEscape([]rune{':', '@'}, credentials[credPos+1:]))
					break
				}
			}
			if !hasSeparator {
				d.User = string(unEscape([]rune{':', '@'}, credentials))
			}

			dsn = dsn[dsnPos+1:]
			break
		}
	}

	// Transport parsing
	for dsnPos, dsnSymbol := range dsn {
		if dsnSymbol != '(' {
			continue
		}

		hpExtractBeginPos := dsnPos + 1
		hpExtractEndPos := -1
		for hpPos, hpSymbol := range dsn[hpExtractBeginPos:] {
			if hpSymbol == ')' {
				hpExtractEndPos = dsnPos + hpPos
			}
		}
		if hpExtractEndPos == -1 {
			continue
		}

		d.Transport = string(dsn[:hpExtractBeginPos-1])
		dsn = append(dsn[hpExtractBeginPos:hpExtractEndPos+1], dsn[hpExtractEndPos+2:]...)
		break
	}

	// Host and port parsing
	for dsnPos, dsnSymbol := range dsn {
		endPos := -1
		if dsnSymbol == '/' {
			endPos = dsnPos
		} else if dsnPos == len(dsn)-1 {
			endPos = len(dsn)
		}

		if endPos > -1 {
			hostPort := dsn[0:endPos]

			hasSeparator := false
			for hpPos, hpSymbol := range hostPort {
				if hpSymbol == ':' {
					hasSeparator = true
					d.Host = string(hostPort[0:hpPos])
					d.Port = string(hostPort[hpPos+1:])
					break
				}
			}
			if !hasSeparator {
				d.Host = string(hostPort)
			}

			dsn = dsn[dsnPos+1:]
			break
		}
	}

	// Path parsing
	for pos, symbol := range dsn {
		endPos := -1
		if symbol == '?' {
			endPos = pos
		} else if pos == len(dsn)-1 {
			endPos = len(dsn)
		}

		if endPos > -1 {
			d.Database = string(dsn[0:endPos])
			dsn = dsn[pos+1:]
			break
		}
	}

	// Params parsing
	beginPosParam := 0
	for symbolPos, symbol := range dsn {
		param := []rune{}
		if symbol == '&' && !isEscaped(symbolPos, dsn) {
			param = dsn[beginPosParam:symbolPos]
			beginPosParam = symbolPos + 1
		} else if symbolPos == len(dsn)-1 {
			param = dsn[beginPosParam:]
		}

		// Separating key and value
		if len(param) > 0 {
			paramKey := []rune{}
			paramVal := []rune{}

			hasSeparator := false
			for paramSymbolPos, paramSymbol := range param {
				if paramSymbol == '=' && !isEscaped(paramSymbolPos, param) {
					hasSeparator = true
					paramKey = param[0:paramSymbolPos]
					paramVal = param[paramSymbolPos+1:]
					break
				}
			}
			if !hasSeparator {
				paramKey = param
			}

			if len(paramKey) > 0 {
				d.AdditionalParams[string(unEscape([]rune{'=', '&'}, paramKey))] = string(unEscape([]rune{'=', '&'}, paramVal))
			}
		}
	}

	return &d
}

func isEscaped(pos int, target []rune) bool {
	return pos > 0 && target[pos-1] == '\\'
}

func searchRuneInArray(target []rune, needle rune) int {
	for i, item := range target {
		if needle == item {
			return i
		}
	}

	return -1
}

func unEscape(needs []rune, target []rune) []rune {
	var unescaped []rune

	for symbolPos, symbol := range target {
		if symbol == '\\' {
			if symbolPos+1 < len(target) {
				if searchRuneInArray(needs, target[symbolPos+1]) != -1 {
					continue
				}
			}
			unescaped = append(unescaped, '\\')
		} else {
			unescaped = append(unescaped, symbol)
		}
	}

	return unescaped
}
