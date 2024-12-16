package core

import (
	"errors"
	"fmt"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/tasks"
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

func (c *Core) GetFlowFromLogID(logID string) (models.Flow, error) {
	f, ok := c.logMap[logID]
	if !ok {
		return models.Flow{}, ErrFlowNotFound
	}

	return c.GetFlowByID(f)
}

// QueueFlowExecution adds a flow in the execution queue. The ID returned is the execution queue ID.
// Log ID should be universally unique, this is used to create the log stream
func (c *Core) QueueFlowExecution(f models.Flow, input map[string]interface{}, logID string) (string, error) {

	// store the mapping between logID and flowID
	c.logMap[logID] = f.Meta.ID

	task, err := tasks.NewFlowExecution(f, input, logID)
	if err != nil {
		return "", fmt.Errorf("error creating task: %v", err)
	}

	info, err := c.q.Enqueue(task)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}
