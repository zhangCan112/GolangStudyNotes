package negroni

// Recovery 是一个可以让程序从任何panic崩溃中恢复的中间件，如果发生panic还会写入一个500错误
type Recovery struct {
}
