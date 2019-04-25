package negroni

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* 测试辅助函数 */
func expect(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestNegroniRun(t *testing.T) {
	go New().Run(":3000")
}

func TestNegroniWith(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n1 := New()
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result = "one"
		next(rw, r)
	}))
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "two"
		next(rw, r)
	}))

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 2, len(n1.Handlers()))
	expect(t, result, "onetwo")

	n2 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "three"
		next(rw, r)
	}))

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 2, len(n1.Handlers()))
	expect(t, result, "onetwo")

	n2.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n2.Handlers()))
	expect(t, result, "onetwothree")
}

func TestNegronWith_doNotModifyOriginal(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n1 := New()
	n1.handlers = make([]Handler, 0, 10) //强行初始化容量
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result = "one"
		next(rw, r)
	}))

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 1, len(n1.Handlers()))

	n2 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "two"
		next(rw, r)
	}))

	n3 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "three"
		next(rw, r)
	}))

	// 重复build 中间件
	n2.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	n3.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 1, len(n1.Handlers()))
	expect(t, result, "one")

	n2.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n2.Handlers()))
	expect(t, result, "onetwo")

	n3.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n3.Handlers()))
	expect(t, result, "onethree")
}
