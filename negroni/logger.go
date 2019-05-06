package negroni

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// LoggerEntry 是传递给模板的结构
type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	HostName  string
	Method    string
	Path      string
	Request   *http.Request
}

// LoggerDefaultDateFormat 是被用作默认的logger 时间格式
var LoggerDefaultDateFormat = time.RFC3339

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger 是一个中间件处理程序，它在请求进入时记录Request，在请求退出时记录Response
type Logger struct {
	// ALogger 实现了足够的log.logger接口，以便与其他实现兼容
	ALogger
	dateFormat string
	template   *template.Template
}

// NewLogger 返回一个新的Logger实例
func NewLogger() *Logger {
	logger := &Logger{ALogger: log.New(os.Stdout, "[negroni]", 0), dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultDateFormat)
	return logger
}

// SetFormat 设置模板格式
func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

// ServeHTTP
func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(ResponseWriter)
	log := LoggerEntry{
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		HostName:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Println(buff.String())
}
