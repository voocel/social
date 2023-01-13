package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"social/pkg/log"
	"time"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(c *gin.Context) {
	now := time.Now()
	path := c.Request.URL.Path
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	method := c.Request.Method
	ip := c.ClientIP()

	buf := new(bytes.Buffer)
	buf.Grow(1024)
	io.Copy(buf, c.Request.Body)
	c.Request.Body = io.NopCloser(buf)
	log.Infow("request",
		log.Pair("request_id", requestID),
		log.Pair("host", ip),
		log.Pair("path", path),
		log.Pair("method", method),
		log.Pair("body", buf.String()),
	)

	bw := &bodyWriter{
		ResponseWriter: c.Writer,
		body:           new(bytes.Buffer),
	}
	c.Writer = bw

	c.Next()

	body := bw.body.String()
	latency := time.Since(now)
	log.Infow("response",
		log.Pair("request_id", requestID),
		log.Pair("host", ip),
		log.Pair("path", path),
		log.Pair("cost", latency),
		log.Pair("body", body),
	)
}
