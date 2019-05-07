package negroni

import (
	"fmt"
	"net/http"
)

const (
	nilRequestMessage = "Request is nil"
)

// PanicInformation 包含用于打印堆栈信息的所有元素
type PanicInformation struct {
	RecoveredPanic interface{}
	Stack          []byte
	Request        *http.Request
}

// StackAsString 返回堆栈的可打印版本
func (p *PanicInformation) StackAsString() string {
	return string(p.Stack)
}

// RequestDescription 返回一个可打印的url
func (p *PanicInformation) RequestDescription() string {
	if p.Request == nil {
		return nilRequestMessage
	}

	var queryOutput string
	if p.Request.URL.RawQuery != "" {
		queryOutput = "?" + p.Request.URL.RawQuery
	}
	return fmt.Sprint("%s %s%s", p.Request.Method, p.Request.URL.Path, queryOutput)
}

// Recovery 是一个可以让程序从任何panic崩溃中恢复的中间件，如果发生panic还会写入一个500错误
type Recovery struct {
	Logger           ALogger
	PrintStack       bool
	LogStack         bool
	PaincHandlerFunc func(*PanicInformation)
	StackAll         bool
	StackSize        int
	Formatter        PanicFormatter
}
