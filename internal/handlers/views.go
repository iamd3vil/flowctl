package handlers

import (
	"io"
	"net/http"
	"html/template"
	"github.com/labstack/echo/v4"
)

type Template struct {
    Templates *template.Template
}

func (t Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}


func (h *Handler) HandleLoginView(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func (h *Handler) HandleFlowsListView(c echo.Context) error {
	flows, err := h.co.GetAllFlows()
	if err != nil {
		h.logger.Error("could not get flows", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "flows", flows)
}
