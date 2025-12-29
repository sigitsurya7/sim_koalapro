package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StreamHandler struct{}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{}
}

func (h *StreamHandler) StatusStream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "stream_unsupported"})
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case t := <-ticker.C:
			payload, _ := json.Marshal(gin.H{"ts": t.UTC().Format(time.RFC3339)})
			_, _ = writer.Write([]byte("event: ping\n"))
			_, _ = writer.Write([]byte("data: "))
			_, _ = writer.Write(payload)
			_, _ = writer.Write([]byte("\n\n"))
			flusher.Flush()
		}
	}
}
