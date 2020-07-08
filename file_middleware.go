package tgin

import (
	"net/http"
	"os"
	"path"
	"strings"
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

func StaticFileMiddleware(urlPrefix, root string, indexes bool) RouteHandler {
	fs := &localFileSystem{
		FileSystem: Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}

	fileServer := http.FileServer(fs)
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
