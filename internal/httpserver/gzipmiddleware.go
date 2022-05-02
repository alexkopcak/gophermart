package httpserver

import (
	"github.com/gin-gonic/gin"
)

func gzipMiddlewareHandle(c *gin.Context) {
	// if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
	// 	gzr, err := gzip.NewReader(c.Request.Body)
	// 	if err != nil {
	// 		c.String(http.StatusBadRequest, "")
	// 		return
	// 	}
	// 	c.Request.Body = gzr
	// }
	c.Next()
}
