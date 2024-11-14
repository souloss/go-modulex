package database

import (
	"fmt"
	"testing"
)

func TestGenDSN(t *testing.T) {
	ds := DataSource{
		Scheme:   "postgres",
		User:     "root",
		Host:     "localhost",
		Password: "root#123",
		Database: "mydb",
	}

	fmt.Println(ds.GenerateDSN())

	ds = DataSource{
		Scheme: "sqlite3",
	}

	fmt.Println(ds.GenerateDSN())
}

func TestParse(t *testing.T) {
	dsns := []string{
		"user:password@example.local",
		"example.local:65535",
		"user:password@",
		"socket:///foo/bar.sock",
		"mysql://user:password@example.local/dbname",
		"mysql://example.local/?db=dbname&user=user&password=password",
	}
	for _, dsn := range dsns {
		ds := Parse(dsn)
		fmt.Println(ds)
	}
}
