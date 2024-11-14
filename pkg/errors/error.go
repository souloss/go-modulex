package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/template"
)

const (
	templateName string = "error"
)

type Error struct {
	error
	// 错误码，解耦具体文字方便实现国际化
	code ErrorCode
	// 格式化字符串以及数据
	message string
	data    map[string]interface{}
	// 堆栈信息
	stack []string
}

func (e *Error) GetCode() ErrorCode {
	return e.code
}

func (e *Error) ToMessage() string {
	var buf bytes.Buffer
	errTemplate, err := template.New(templateName).Parse(e.message)
	if err != nil {
		panic(err)
	}
	err = errTemplate.Execute(&buf, e.data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (e *Error) Error() string {
	var sb strings.Builder
	var buf bytes.Buffer
	errTemplate, err := template.New(templateName).Parse(e.message)
	if err != nil {
		panic(err)
	}
	err = errTemplate.Execute(&buf, e.data)
	if err != nil {
		panic(err)
	}
	sb.WriteString(fmt.Sprintf("Error Code: %d\nMessage: %s\n", e.code, buf.String()))
	if len(e.stack) > 0 {
		sb.WriteString("Stack Trace:\n")
		for _, frame := range e.stack {
			sb.WriteString(fmt.Sprintf("\t%s\n", frame))
		}
	}
	return sb.String()
}

// New 创建一个新的错误，并附带业务错误码
func New(code ErrorCode) *Error {
	return &Error{
		code:    code,
		message: code.String(),
		data:    make(map[string]interface{}),
		stack:   captureStack(),
	}
}

// New 创建一个新的错误，并附带业务错误码
func NewWithData(code ErrorCode, data map[string]interface{}) *Error {
	return &Error{
		code:    code,
		message: code.String(),
		data:    data,
		stack:   captureStack(),
	}
}

// Wrap 包装一个已有的错误
func Wrap(err error) *Error {
	return &Error{
		code:    ErrUndefined,
		message: ErrUndefined.String(),
		data: map[string]interface{}{
			"error": err.Error(),
		},
		stack: captureStack(),
	}
}

func captureStack() []string {
	const size = 32
	var pcs [size]uintptr           // 存储调用栈信息
	n := runtime.Callers(3, pcs[:]) // 跳过前三层的堆栈 (runtime.Callers, captureStack, 生成错误函数)
	stack := make([]string, 0, n)

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		stack = append(stack, fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return stack
}
