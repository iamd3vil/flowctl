package runner

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// StreamLogger logs messages to a Redis stream
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

func (s *StreamLogger) Close() error {
	return s.r.XAdd(context.Background(), &redis.XAddArgs{
		Stream: s.ID,
		Values: map[string]interface{}{"closed": true},
	}).Err()
}
