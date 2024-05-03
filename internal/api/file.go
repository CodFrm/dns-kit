package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/dns-kit/frontend"
	"github.com/gin-gonic/gin"
)

// File 处理静态资源
func File(root *mux.Router) error {
	root.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}
		path := c.Request.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		data, err := frontend.EmbedFS.ReadFile("dist" + path)
		if err != nil {
			c.Next()
			if c.Writer.Status() != http.StatusNotFound {
				return
			}
			// 重写到index.html
			data, err = frontend.EmbedFS.ReadFile("dist/index.html")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}
		c.Status(http.StatusOK)
		switch filepath.Ext(path) {
		case ".html":
			c.Header("Content-Type", "text/html")
		case ".js":
			c.Header("Content-Type", "application/javascript")
		case ".css":
			c.Header("Content-Type", "text/css")
		default:
			c.Header("Content-Type", http.DetectContentType(data))
		}
		_, _ = c.Writer.Write(data)
	})
	return nil
}
