package negroni

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"text/template"
)

const (
	// NoPrintStackBodyString is the body content returned when HTTP stack printing is suppressed
	NoPrintStackBodyString = "500 Internal Server Error"
	panicText              = "PANIC: %s\n%s"
	nilRequestMessage      = "Request is nil"
	panicHTML              = `<html>
<head><title>PANIC: {{.RecoveredPanic}}</title></head>
<style type="text/css">
html, body {
	font-family: Helvetica, Arial, Sans;
	color: #333333;
	background-color: #ffffff;
	margin: 0px;
}
h1 {
	color: #ffffff;
	background-color: #f14c4c;
	padding: 20px;
	border-bottom: 1px solid #2b3848;
}
.block {
	margin: 2em;
}
.panic-interface {
}

.panic-stack-raw pre {
	padding: 1em;
	background: #f6f8fa;
	border: dashed 1px;
}
.panic-interface-title {
	font-weight: bold;
}
</style>
<body>
<h1>Negroni - PANIC</h1>

<div class="panic-interface block">
	<h3>{{.RequestDescription}}</h3>
	<span class="panic-interface-title">Runtime error:</span> <span class="panic-interface-element">{{.RecoveredPanic}}</span>
</div>

{{ if .Stack }}
<div class="panic-stack-raw block">
	<h3>Runtime Stack</h3>
	<pre>{{.StackAsString}}</pre>
</div>
{{ end }}

</body>
</html>`
)

var panicHTMLTemplate = template.Must(template.New("PanicPage").Parse(panicHTML))

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

// PanicFormatter 是对象上的接口，可以实现用来输出堆栈的跟踪信息
type PanicFormatter interface {
	// FormatPanicError 为给定的应答/响应输出堆栈
	// 如果中间件不输出堆栈跟踪信息
	// 那么传递的PanicInformation 的Stack 为空的字节数组[]btye{}
	FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation)
}

// TextPanicFormatter 在os.stdout上将堆栈输出为简单文本。
// 如果未设置“content type”，则将数据输出为“text/plain；charset=utf-8”。
// 否则，将保留源代码“content type”。
type TextPanicFormatter struct{}

// FormatPanicError 实现PanicFormatter接口方法
func (t *TextPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) {
	if rw.Header().Get("Content-type") == "" {
		rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	}
	fmt.Fprintf(rw, panicText, infos.RecoveredPanic, infos.Stack)
}

// HTMLPanicFormatter 输出堆栈信息到HTML页面内。
// 这在很大程度上受到了
// https://github.com/go-martini/martini/pull/156/commits的启发。
type HTMLPanicFormatter struct{}

// FormatPanicError 实现PanicFormatter接口方法
func (t *HTMLPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) {
	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	panicHTMLTemplate.Execute(rw, infos)
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

	// Deprecated: 请改用PanicHandlerFunc
	// 接收包含附加信息的panic错误(请参阅PanicInformation)
	ErrorHandleFunc func(interface{})
}

// NewRecovery 返回一个新的Recovery实例
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[negroni]", 0),
		PrintStack: true,
		LogStack:   true,
		StackAll:   false,
		StackSize:  1024 * 8,
		Formatter:  &TextPanicFormatter{},
	}
}

// ServeHttp Recovery的http.Handler接口实现
func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			//他认为他给的Size足够大，才这么操作的
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			infos := &PanicInformation{RecoveredPanic: err, Request: r}

			if rec.PrintStack {
				infos.Stack = stack
				rec.Formatter.FormatPanicError(rw, r, infos)
			} else {
				if rw.Header().Get("Content-Type") == "" {
					rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}
				fmt.Fprint(rw, NoPrintStackBodyString)
			}

			if rec.LogStack {
				rec.Logger.Printf(panicText, err, stack)
			}

			if rec.ErrorHandleFunc != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							rec.Logger.Printf("provided ErrorHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack())
							rec.Logger.Printf("%s\n", debug.Stack())
						}
					}()
					rec.ErrorHandleFunc(err)
				}()
			}

			if rec.PaincHandlerFunc != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							rec.Logger.Printf("provided PanicHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack())
							rec.Logger.Printf("%s\n", debug.Stack())
						}
					}()
					rec.PaincHandlerFunc(infos)
				}()
			}
		}
	}()

	next(rw, r)
}
