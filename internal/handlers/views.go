package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cvhariharan/autopilot/internal/core"
	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/cvhariharan/autopilot/internal/ui/partials"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
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

func (h *Handler) HandleFlowTrigger(c echo.Context) error {
	var req map[string]interface{}
	// This is done to only bind request body and ignore path / query params
	if err := (&echo.DefaultBinder{}).BindBody(c, &req); err != nil {
		return showErrorPage(c, http.StatusBadRequest, "could not parse request")
	}

	f, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		return showErrorPage(c, http.StatusNotFound, err.Error())
	}

	if err := f.ValidateInput(req); err != nil {
		return render(c, ui.FlowInputFormPage(f, "", map[string]string{err.FieldName: err.Msg}, ""))
	}

	// Add to queue
	logID := uuid.NewString()
	_, err = h.co.QueueFlowExecution(f, req, logID)
	if err != nil {
		return render(c, ui.FlowInputFormPage(f, "", nil, err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/view/results/%s/%s", f.Meta.ID, logID))
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) HandleFlowForm(c echo.Context) error {
	flow, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		return showErrorPage(c, http.StatusNotFound, err.Error())
	}

	return ui.FlowInputFormPage(flow, "", nil, "").Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) HandleFlowsList(c echo.Context) error {
	flows, err := h.co.GetAllFlows()
	if err != nil {
		return showErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	return render(c, ui.FlowsListPage(flows))
}

func (h *Handler) HandleFlowExecutionResults(c echo.Context) error {
	flowID := c.Param("flowID")
	if flowID == "" {
		return showErrorPage(c, http.StatusBadRequest, "flow id cannot be empty")
	}

	logID := c.Param("logID")
	if logID == "" {
		return showErrorPage(c, http.StatusBadRequest, "execution id cannot be empty")
	}

	f, err := h.co.GetFlowByID(flowID)
	if err != nil {
		return showErrorPage(c, http.StatusNotFound, err.Error())
	}

	return render(c, ui.ResultsPage(f.Meta.Name, fmt.Sprintf("/view/logs/%s", logID)))
}

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

	msgCh := h.co.StreamLogs(c.Request().Context(), c.Param("logID"))

	for msg := range msgCh {
		if msg.Err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := partials.LogMessage(msg.Message).Render(c.Request().Context(), &buf); err != nil {
			return err
		}

		if err := ws.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
