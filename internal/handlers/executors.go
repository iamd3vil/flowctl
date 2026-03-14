package handlers

import (
	"net/http"

	"github.com/cvhariharan/flowctl/sdk/executor"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleGetExecutorConfig(c echo.Context) error {
	executorName := c.Param("executor")
	if executorName == "" {
		return wrapError(ErrRequiredFieldMissing, "executor name cannot be empty", nil, nil)
	}

	schema, err := executor.GetSchema(executorName)
	if err != nil {
		return wrapError(ErrResourceNotFound, "could not get executor config", err, nil)
	}

	return c.JSON(http.StatusOK, schema)
}

func (h *Handler) HandleListExecutors(c echo.Context) error {
	entries := executor.GetAllExecutors()
	infos := make([]ExecutorInfo, 0, len(entries))
	for _, e := range entries {
		infos = append(infos, ExecutorInfo{
			Name:         e.Name,
			Capabilities: e.Capabilities,
		})
	}
	return c.JSON(http.StatusOK, ExecutorsListResponse{Executors: infos})
}
