package core

import (
	"errors"

	"github.com/cvhariharan/autopilot/internal/models"
)

var (
	ErrFlowNotFound = errors.New("flow not found")
)

func (c *Core) GetFlowByID(id string) (models.Flow, error) {
	f, ok := c.flows[id]
	if !ok {
		return models.Flow{}, ErrFlowNotFound
	}

	return f, nil
}

func (c *Core) GetAllFlows() ([]models.Flow, error) {
	var fs []models.Flow
	for _, v := range c.flows {
		fs = append(fs, v)
	}
	return fs, nil
}
