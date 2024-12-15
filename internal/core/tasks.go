package core

import (
	"fmt"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/tasks"
)

// QueueFlowExecution adds a flow in the execution queue. The ID returned is the execution queue ID.
// Log ID should be universally unique, this is used to create the log stream
func (c *Core) QueueFlowExecution(f models.Flow, input map[string]interface{}, logID string) (string, error) {
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
