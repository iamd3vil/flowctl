package handlers

import (
	"net/http"

	"github.com/cvhariharan/flowctl/internal/messengers"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleGetMessengers(c echo.Context) error {
	return c.JSON(http.StatusOK, messengers.GetAllSchemas())
}
