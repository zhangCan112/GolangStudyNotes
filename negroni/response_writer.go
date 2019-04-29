package negroni

import "net/http"

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
	Before(fuc ResponseWriter)
}

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

type responseWriterCloseNotifer struct {
	*responseWriter
}

func (rw *responseWriterCloseNotifer) closeNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
