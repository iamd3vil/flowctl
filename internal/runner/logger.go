package runner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/redis/go-redis/v9"
)

const CheckpointPrefix = "checkpoint:%s"

// StreamLogger logs messages to a Redis stream and provides utility functions for tracking execution status
type StreamLogger struct {
	ID string
	r  redis.UniversalClient
}

// NewStreamLogger creates a new StreamLogger instance with the given Redis client
func NewStreamLogger(r redis.UniversalClient) *StreamLogger {
	return &StreamLogger{r: r}
}

// WithID sets the ID for the logger
// The ID is used to identify the stream. Can be used to separate actions
func (s *StreamLogger) WithID(id string) *StreamLogger {
	s.ID = id
	return s
}

func (s *StreamLogger) Write(p []byte) (int, error) {
	res := s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"log": string(p)},
	})
	if res.Err() != nil {
		return 0, res.Err()
	}
	return len(p), nil
}

// Checkpoint can be used to save the completion status of an action.
// Call after the successful completion of an action
func (s *StreamLogger) Checkpoint(id string, chck models.ExecutionCheckpoint) error {
	chck.ActionID = id
	if chck.Results != nil {
		res, err := json.Marshal(chck.Results)
		if err != nil {
			return fmt.Errorf("could not marshal checkpoint results for %s: %w", id, err)
		}
		if _, err := s.r.Set(context.Background(), fmt.Sprintf(CheckpointPrefix, id), res, 0).Result(); err != nil {
			return fmt.Errorf("error while storing checkpoint results for %s: %w", id, err)
		}
	}
	return s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"checkpoint": chck},
	}).Err()
}

func (s *StreamLogger) Results(id string) (map[string]string, error) {
	resB, err := s.r.Get(context.Background(), id).Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not get results for %s: %w", id, err)
	}

	var results map[string]string
	if err := json.Unmarshal(resB, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *StreamLogger) Close() error {
	return s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"closed": true},
	}).Err()
}
