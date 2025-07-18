package streamlogger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const CheckpointPrefix = "checkpoint:%s"

// StreamLogger logs messages to a Redis stream and provides utility functions for tracking execution status
// This should only be used when actions are executed sequentially, this is NOT thread-safe
type StreamLogger struct {
	ID string
	r  redis.UniversalClient
	actionID string
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

// WithActionID sets the current action ID.
// This is used to tag messages. This is NOT thread-safe
func (s *StreamLogger) WithActionID(id string) *StreamLogger {
	s.actionID = id
	return s
}

// Write to a redis stream, this implementation is NOT thread-safe.
// Write uses the current actionID to mark the messages
func (s *StreamLogger) Write(p []byte) (int, error) {
	if err := s.Checkpoint(s.actionID, p, LogMessageType); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Checkpoint can be used to save the completion status of an action.
// Call after the successful completion of an action
func (s *StreamLogger) Checkpoint(id string, val interface{}, mtype MessageType) error {
	var sm StreamMessage
	sm.ActionID = id
	switch mtype {
	case ErrMessageType:
		e, ok := val.(string)
		if !ok {
			return fmt.Errorf("expected string type for error got %T in stream checkpoint", val)
		}
		sm.MType = ErrMessageType
		sm.Val = []byte(e)
	case ResultMessageType:
		r, ok := val.(map[string]string)
		if !ok {
			return fmt.Errorf("expected map[string]string type got %T in stream checkpoint", val)
		}
		data, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("could not marshal result for result message type in stream message %s: %w", id, err)
		}
		sm.MType = ResultMessageType
		sm.Val = data
	case LogMessageType:
		sm.MType = LogMessageType
		d, ok := val.([]byte)
		if !ok {
			return fmt.Errorf("expected []byte type for log got %T in stream checkpoint", val)
		}
		sm.MType = LogMessageType
		sm.Val = d
	}

	return s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"checkpoint": sm},
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

func (s *StreamLogger) Close(closeID string) error {
	return s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"closed": closeID},
	}).Err()
}
