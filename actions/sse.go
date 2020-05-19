package actions

import (
	. "culture/sse"
	"errors"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//SseHandler SseHandler
func SseHandler(c buffalo.Context) error {
	// We need to be able to flush for SSE
	fl, ok := c.Response().(http.Flusher)
	if !ok {
		http.Error(c.Response(), "Flushing not supported", http.StatusNotImplemented)
		return c.Error(http.StatusNotImplemented, errors.New(http.StatusText(http.StatusNotImplemented)))
	}

	// Returns a channel that blocks until the connection is closed
	done := c.Request().Context().Done()

	// Set headers for SSE
	h := c.Response().Header()
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")

	// Connect new client
	cl := make(SseClient, S.BufSize)
	S.Connecting <- cl

	for {
		select {
		case <-done:
			// Disconnect the client when the connection is closed
			S.Disconnecting <- cl
			return nil

		case event := <-cl:
			// Write events
			c.Response().Write(event)
			fl.Flush()
		}
	}
}

//UpdateSystem UpdateSystem
func UpdateSystem(c buffalo.Context) error {
	S.SendString("", "update_system", "update_system")
	return c.Render(http.StatusOK, nil)
}
