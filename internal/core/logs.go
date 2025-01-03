package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/redis/go-redis/v9"
)

// StreamLogs reads values from a redis stream from the beginning and returns a channel to which
// all the messages are sent. logID is the ID sent to the NewFlowExecution task
func (c *Core) StreamLogs(ctx context.Context, logID string) chan models.StreamMessage {
	ch := make(chan models.StreamMessage)

	go func(ch chan models.StreamMessage) {
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
				ch <- models.StreamMessage{MType: models.ErrMessageType, Val: []byte(fmt.Errorf("error reading from stream: %w", err).Error())}
				return
			}

			for _, stream := range result {
				for _, message := range stream.Messages {
					if _, ok := message.Values["closed"]; ok {
						return
					}

					if checkpoint, ok := message.Values["checkpoint"]; ok {
						sm, ok := checkpoint.(models.StreamMessage)
						if !ok {
							log.Printf("checkpoint not of StreamMessage type: %v", checkpoint)
							continue
						}
						ch <- sm
					}

					lastProcessedID = message.ID
				}
			}

			time.Sleep(1 * time.Second)
		}

	}(ch)

	return ch
}
