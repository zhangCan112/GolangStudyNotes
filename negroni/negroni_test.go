package negroni

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestNegroniServeHttp(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n := New()
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "foo"
		next(rw, r)
		result += "ban"
	}))

	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bar"
		next(rw, r)
		result += "baz"
	}))

	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bat"
		rw.WriteHeader(http.StatusBadRequest)
	}))

	n.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

// 确保Negroni中间件链可以正确的返回他的所有handler
func TestHandlers(t *testing.T) {
	response := httptest.NewRecorder()
	n := New()
	handlers := n.Handlers()
	expect(t, 0, len(handlers))

	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.WriteHeader(http.StatusOK)
	}))

	// 期望在添加了一个handler之后，handlers能返回长度为1
	handlers = n.Handlers()
	expect(t, 1, len(handlers))

	// 确保早注册的hander在handlers数组中也是前面的
	handlers[0].ServeHTTP(response, (*http.Request)(nil), nil)
	expect(t, response.Code, http.StatusOK)
}

func TestDetectAddress(t *testing.T) {
	if detectAddress() != DefaultAddress {
		t.Error("Expected the DefaultAddress")
	}

	if detectAddress(":6060") != ":6060" {
		t.Error("Expected the provided address")
	}

	os.Setenv("PORT", "9090")
	if detectAddress() != ":9090" {
		t.Error("Expected the PORT env var with a prefixed colon")
	}
}

func voidHTTPHandlerFunc(rw http.ResponseWriter, r *http.Request) {
	// Do nothing
}

// 测试Wrap函数
func TestWrap(t *testing.T) {
	reponse := httptest.NewRecorder()

	handler := Wrap(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(reponse, (*http.Request)(nil), voidHTTPHandlerFunc)

	expect(t, reponse.Code, http.StatusOK)

}

// 测试WrapFunc函数
func TestWrapFunc(t *testing.T) {
	response := httptest.NewRecorder()

	handler := WrapFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(response, (*http.Request)(nil), voidHTTPHandlerFunc)

	expect(t, response.Code, http.StatusOK)
}
