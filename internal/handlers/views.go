package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/core"
	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/cvhariharan/autopilot/internal/ui/partials"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	flows map[string]models.Flow
	co    *core.Core
}

func NewHandler(f map[string]models.Flow, co *core.Core) *Handler {
	return &Handler{flows: f, co: co}
}

func (h *Handler) HandleTrigger(c echo.Context) error {
	var req map[string]interface{}
	// This is done to only bind request body and ignore path / query params
	if err := (&echo.DefaultBinder{}).BindBody(c, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error validating request bind")
	}

	f, ok := h.flows[c.Param("flow")]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "requested flow not found")
	}

	if err := f.ValidateInput(req); err != nil {
		var ferr *models.FlowValidationError
		if errors.As(err, &ferr) {
			return ui.FlowInputForm(f, map[string]string{ferr.FieldName: ferr.Msg}).Render(c.Request().Context(), c.Response().Writer)
		}
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error validating input: %v", err))
	}

	// Add to queue
	logID := uuid.NewString()
	_, err := h.co.QueueFlowExecution(f, req, logID)
	if err != nil {
		return err
	}

	return partials.LogTerminal(fmt.Sprintf("/api/logs/%s", logID)).Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) HandleForm(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error validating request bind")
	}
	flow, ok := h.flows[c.Param("flow")]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "requested flow not found")
	}

	return ui.FlowInputForm(flow, make(map[string]string)).Render(c.Request().Context(), c.Response().Writer)
}
