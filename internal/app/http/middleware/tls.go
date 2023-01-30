package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"social/pkg/log"
)

func TLSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := s.Process(c.Writer, c.Request)
		if err != nil {
			log.Error(err)
			return
		}
		c.Next()
	}
}
