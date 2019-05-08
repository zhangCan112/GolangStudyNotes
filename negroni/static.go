package negroni

import (
	"net/http"
	"path"
	"strings"
)

// Static 是为给定目录/文件系统中的静态文件提供服务的中间件处理程序。如果文件系统上不存在该文件，则
// 传递到链中的下一个中间件。如果你想要“文件服务器”
// 当它为未找到的文件返回404时，您应该考虑
// 使用go stdlib中的http.fileserver。
type Static struct {
	// Dir 是为静态文件提供服务的目录
	Dir http.FileSystem
	// Prefix 是用于服务静态目录内容的可选前缀
	Prefix string
	// IndexFile 定义要用作索引的文件（如果存在）。
	IndexFile string
}

// NewStatic 返回一个新的 Static实例
func NewStatic(directory http.FileSystem) *Static {
	return &Static{
		Dir:       directory,
		Prefix:    "",
		IndexFile: "index.html",
	}
}

func (s *Static) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}
	file := r.URL.Path

	if s.Prefix != "" {
		if !strings.HasPrefix(file, s.Prefix) {
			next(rw, r)
			return
		}
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' {
			next(rw, r)
			return
		}
	}

	f, err := s.Dir.Open(file)
	if err != nil {
		// discard the error?
		next(rw, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}
	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil {
			next(rw, r)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(rw, r)
			return
		}
	}

	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
