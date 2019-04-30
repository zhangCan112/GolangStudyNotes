package negroni

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// ResponseWriter 是一个围绕http.ResponseWriter 提供额外信息的包装器。
// 如果有需要的话推荐中间件用这个结构去包装ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status 返回response的status码或者0(当response还未写入时)
	Status() int
	// Written 返回ResponseWriter是否已被写入过
	Written() bool
	// Size 返回response body的大小
	Size() int
	// Before 允许在写入ResponseWriter之前调用函数，
	// 这对于必须在Response的写操作之前设置Header或者其他操作很有用
	Before(func(ResponseWriter))
}

// NewResponseWriter 包装http.ResponseWriter来创建一个ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	if _, ok := rw.(http.CloseNotifier); ok {
		return &responseWriterCloseNotifer{nrw}
	}

	return nrw
}

type beforeFunc func(ResponseWriter)
type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, error := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, error
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	fluser, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			rw.WriteHeader(http.StatusOK)
		}
		fluser.Flush()
	}
}

type responseWriterCloseNotifer struct {
	*responseWriter
}

func (rw *responseWriterCloseNotifer) closeNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
