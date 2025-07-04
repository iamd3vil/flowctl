package handlers

import (
	"fmt"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/core/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleFlowTrigger(c echo.Context) error {
	user, ok := c.Get("user").(models.UserInfo)
	if !ok {
		return wrapError(http.StatusBadRequest, "could not get user details", nil, nil)
	}

	var req map[string]interface{}
	// This is done to only bind request body and ignore path / query params
	if err := (&echo.DefaultBinder{}).BindBody(c, &req); err != nil {
		return wrapError(http.StatusBadRequest, "could not parse request", err, nil)
	}

	f, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		return wrapError(http.StatusBadRequest, "could not get flow", err, nil)
	}

	if len(f.Actions) == 0 {
		return wrapError(http.StatusBadRequest, "no actions in flow", nil, nil)
	}

	if err := f.ValidateInput(req); err != nil {
		return wrapError(http.StatusBadRequest, "", err, FlowInputValidationError{
			FieldName: err.FieldName,
			ErrMessage: err.Msg,
		})
	}

	// Add to queue
	execID, err := h.co.QueueFlowExecution(c.Request().Context(), f, req, user.ID)
	if err != nil {
		return wrapError(http.StatusBadRequest, "could not trigger flow", err, nil)
	}

	if err := c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/view/results/%s/%s", f.Meta.ID, execID)); err != nil {
		h.logger.Error("redirect", "error", err)
		return err
	}
	return c.NoContent(http.StatusOK)
}
