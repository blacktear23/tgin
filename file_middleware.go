package tgin

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

const INDEX = "index.html"

type noDirFs struct {
	fs http.FileSystem
}

type noDirFile struct {
	http.File
}

func (fs *noDirFs) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &noDirFile{f}, nil
}

func (f *noDirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func Dir(root string, listDir bool) http.FileSystem {
	fs := http.Dir(root)
	if listDir {
		return fs
	}
	return &noDirFs{fs}
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func (l *localFileSystem) Exists(prefix, filePath string) bool {
	if p := strings.TrimPrefix(filePath, prefix); len(p) < len(filePath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(ioutil.Discard)
		gzip.NewWriterLevel(w, gzip.DefaultCompression)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func StaticFileMiddleware(urlPrefix, root string, indexes bool) RouteHandler {
	fs := &localFileSystem{
		FileSystem: Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}

	fileServer := GzipHandler(http.FileServer(fs))
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
