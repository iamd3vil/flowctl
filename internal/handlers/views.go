package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cvhariharan/autopilot/internal/core"
	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/cvhariharan/autopilot/internal/ui/partials"
	"github.com/google/uuid"
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

func (h *Handler) HandleFlowTrigger(c echo.Context) error {
	var req map[string]interface{}
	// This is done to only bind request body and ignore path / query params
	if err := (&echo.DefaultBinder{}).BindBody(c, &req); err != nil {
		return showErrorPage(c, http.StatusNotFound, "could not parse request")
	}

	f, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		return render(c, ui.FlowInputFormPage(f, nil, err.Error()))
	}

	if err := f.ValidateInput(req); err != nil {
		return render(c, ui.FlowInputFormPage(f, map[string]string{err.FieldName: err.Msg}, ""))
	}

	// Add to queue
	logID := uuid.NewString()
	_, err = h.co.QueueFlowExecution(f, req, logID)
	if err != nil {
		return render(c, ui.FlowInputFormPage(f, nil, err.Error()))
	}

	return render(c, partials.LogTerminal(fmt.Sprintf("/api/logs/%s", logID)))
}

func (h *Handler) HandleFlowForm(c echo.Context) error {
	flow, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		return showErrorPage(c, http.StatusNotFound, err.Error())
	}

	return ui.FlowInputFormPage(flow, nil, "").Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) HandleFlowsList(c echo.Context) error {
	flows, err := h.co.GetAllFlows()
	if err != nil {
		return showErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	return render(c, ui.FlowsListPage(flows))
}
