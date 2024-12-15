package handlers

import (
	"bytes"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

// HandleLogStreaming uses SSE to stream action logs
func (h *Handler) HandleLogStreaming(c echo.Context) error {
	// Upgrade to WebSocket connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(int(http.StateClosed), "Connection closed"))
		ws.Close()
	}()

	msgCh := h.core.StreamLogs(c.Request().Context(), c.Param("flow"))

	for msg := range msgCh {
		if msg.Err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := ui.LogMessage(msg.Message).Render(c.Request().Context(), &buf); err != nil {
			return err
		}

		if err := ws.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
