package core

import (
	"context"
	"fmt"
	"time"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/redis/go-redis/v9"
)

// StreamLogs reads values from a redis stream from the beginning and returns a channel to which
// all the messages are sent. logID is the ID sent to the NewFlowExecution task
func (c *Core) StreamLogs(ctx context.Context, logID string) chan models.LogMessage {
	ch := make(chan models.LogMessage)

	go func(ch chan models.LogMessage) {
		defer close(ch)
		lastProcessedID := "0"
		for {
			result, err := c.redisClient.XRead(ctx, &redis.XReadArgs{
				Streams: []string{logID, lastProcessedID},
				Count:   10,
				Block:   0,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				ch <- models.LogMessage{Err: fmt.Sprintf("error reading from redis log stream: %v", err)}
				return
			}

			for _, stream := range result {
				for _, message := range stream.Messages {
					if _, ok := message.Values["closed"]; ok {
						return
					}

					if checkpoint, ok := message.Values["checkpoint"]; ok {
						var chck models.ExecutionCheckpoint
						chck.UnmarshalBinary([]byte(checkpoint.(string)))

						if chck.Err != "" {
							ch <- models.LogMessage{ID: chck.ActionID, Checkpoint: true, Err: chck.Err}
							continue
						}

						ch <- models.LogMessage{ID: chck.ActionID, Checkpoint: true, Results: chck.Results}
						continue
					}

					ch <- models.LogMessage{Message: message.Values["log"].(string)}

					lastProcessedID = message.ID
				}
			}

			time.Sleep(1 * time.Second)
		}

	}(ch)

	return ch
}
