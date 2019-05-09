package context

import (
	"net/http"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	data  = make(map[*http.Request]map[interface{}]interface{})
	datat = make(map[*http.Request]int64)
)

// Set 保存一个给定key的值在一个给定的request里
func Set(r *http.Request, key, val interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if data[r] == nil {
		data[r] = make(map[interface{}]interface{})
		datat[r] = time.Now().Unix()
	}
	data[r][key] = val
}

// Get 返回一个给定request里给定key的值
func Get(r *http.Request, key interface{}) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()

	if ctx := data[r]; ctx != nil {
		value := ctx[key]
		return value
	}
	return nil
}

// GetOK  返回一个指定的值，并且返回该值是否获取成功（是不存在key，还是value为nil）
func GetOK(r *http.Request, key interface{}) (val interface{}, ok bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	if _, ok := data[r]; ok {
		value, ok := data[r][key]
		return value, ok
	}

	return nil, false
}

// GetAll 返回request所有键值对
func GetAll(r *http.Request) map[interface{}]interface{} {
	mutex.RLock()
	defer mutex.RUnlock()

	if context, ok := data[r]; ok {
		result := make(map[interface{}]interface{}, len(context))
		for k, v := range context {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetAllOk 返回request所有键值对,并返回状态表示request是否是已注册
func GetAllOk(r *http.Request) (map[interface{}]interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	context, ok := data[r]
	result := make(map[interface{}]interface{}, len(context))
	for k, v := range context {
		result[k] = v
	}
	return result, ok
}

// Delete 删除指定request下保存的某个指定key的值
func Delete(r *http.Request, key interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if data[r] != nil {
		delete(data[r], key)
	}
}

// Clear 清楚所有的数据
func Clear(r *http.Request) {
	mutex.Lock()
	clear(r)
	mutex.Unlock()
}

func clear(r *http.Request) {
	delete(data, r)
	delete(datat, r)
}

// Purge 根据时间清理保存的context
func Purge(maxAge int) int {
	mutex.Lock()
	defer mutex.Unlock()

	count := 0
	if maxAge <= 0 {
		count = len(data)
		data = make(map[*http.Request]map[interface{}]interface{})
		datat = make(map[*http.Request]int64)
	} else {
		min := time.Now().Unix() - int64(maxAge)
		for r := range data {
			if datat[r] < min {
				clear(r)
				count++
			}
		}
	}
	return count
}

// ClearHandler 包装http.Handler并在请求生存期结束时清除请求值。
func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer Clear(r)
		h.ServeHTTP(rw, r)
	})
}
