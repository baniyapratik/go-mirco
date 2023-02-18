package sever

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) newHTTPHandler() http.Handler {
	r := gin.New()
	// Middlewares
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	return r
}
