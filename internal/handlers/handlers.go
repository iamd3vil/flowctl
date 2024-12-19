package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cvhariharan/autopilot/internal/core"
	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	co *core.Core
}

func NewHandler(co *core.Core) *Handler {
	return &Handler{co: co}
}

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func showErrorPage(c echo.Context, code int, message string) error {
	return ui.ErrorPage(code, message).Render(c.Request().Context(), c.Response().Writer)
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	errMsg := "error processing the request"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		errMsg = he.Message.(string)
	}

	c.Logger().Error(err)

	if err := showErrorPage(c, code, errMsg); err != nil {
		c.Logger().Error(err)
	}
}

func renderToWebsocket(c echo.Context, component templ.Component, ws *websocket.Conn) error {
	var buf bytes.Buffer
	if err := component.Render(c.Request().Context(), &buf); err != nil {
		return fmt.Errorf("could not render component: %w", err)
	}

	if err := ws.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
		return fmt.Errorf("could not send to websocket: %w", err)
	}

	return nil
}
