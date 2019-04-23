package negroni

import "net/http"

// Handler 中间件定义的hanlder接口，比http.Handler多了一个next参数
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// HandlerFunc 就是一个允许普通函数做为handler的适配器，
// 因为将对函数类型添加了方法，所以同签名的函数就可以方面的传入类型为hanler的参数中了
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type middleware struct {
	handler Handler
	next    *middleware
}

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}

// Wrap 用来将http.Handler包装成negroni的Handler
func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

// WrapFunc 用来将http.HandlerFunc包装成negroni的Handler
func WrapFunc(handlerFunc http.HandlerFunc) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handlerFunc.ServeHTTP(rw, r)
		next(rw, r)
	})
}

// Negroni 是一组中间件的处理程序， 可以作为http.handler调用
// negroni中间件按添加到堆栈的顺序进行计算
type Negroni struct {
	middleware middleware
	handlers   []Handler
}

func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(rw, r)
}

// New 返回一个预先没有配置中间件的新的Negroni实例
func New(handlers ...Handler) *Negroni {
	return &Negroni{
		handlers:   handlers,
		middleware: build(handlers),
	}
}

// With 根据当前的Negroni实例中的数据和新的handlers返回一个新的Negroni实例
func (n *Negroni) With(handlers ...Handler) *Negroni {
	currentHandlers := make([]Handler, len(n.handlers))
	copy(currentHandlers, n.handlers)
	return New(
		append(currentHandlers, handlers...)...,
	)
}

// Use 添加一个handler到中间件栈中
func (n *Negroni) Use(handler Handler) {

}

// build 利用递归构建中间件链
func build(handlers []Handler) middleware {
	var next middleware
	switch {
	case len(handlers) == 0:
		return voidMiddleware()
	case len(handlers) > 1:
		next = build(handlers[1:])
	default:
		next = voidMiddleware()
	}
	return middleware{
		handlers[0],
		&next,
	}
}

func voidMiddleware() middleware {
	return middleware{
		handler: HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		next:    &middleware{},
	}
}
