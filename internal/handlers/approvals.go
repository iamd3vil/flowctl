package handlers

import (
	"fmt"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/ui"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleApprovalRequest(c echo.Context) error {
	approvalID := c.Param("approvalID")
	if approvalID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "approval ID cannot be empty")
	}

	areq, err := h.co.GetApprovalRequest(c.Request().Context(), approvalID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get approval request")
	}

	f, err := h.co.GetFlowFromLogID(areq.ExecID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get flow")
	}

	return render(c, ui.ApprovalPage(ui.ApprovalRequest{
		ID:          approvalID,
		Title:       f.Meta.Name,
		Description: fmt.Sprintf("Approval request to execute action %q from flow %q", areq.ActionID, f.Meta.Name),
		RequestedBy: areq.RequestedBy,
	}), http.StatusOK)
}
