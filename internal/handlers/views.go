package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/cvhariharan/autopilot/internal/core/models"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

type Page struct {
	Title      string
	ErrMessage string
	Message    string
}

func (h *Handler) HandleLoginView(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func (h *Handler) HandleFlowsListView(c echo.Context) error {
	data := struct {
		Page
		Flows []models.Flow
	}{
		Page: Page{
			Title: "Flows",
		},
	}

	flows, err := h.co.GetAllFlows()
	if err != nil {
		data.ErrMessage = err.Error()
		return c.Render(http.StatusBadRequest, "flows", data)
	}

	data.Flows = flows
	return c.Render(http.StatusOK, "flows", data)
}

func (h *Handler) HandleFlowFormView(c echo.Context) error {
	data := struct {
		Page
		Flow       models.Flow
		InputErrors map[string]string
	}{
		Page: Page{
			Title: "Flow Input",
		},
	}
	flow, err := h.co.GetFlowByID(c.Param("flow"))
	if err != nil {
		data.ErrMessage = err.Error()
		return c.Render(http.StatusBadRequest, "flow_input", data)
	}
	data.Flow = flow

	return c.Render(http.StatusOK, "flow_input", data)
}
