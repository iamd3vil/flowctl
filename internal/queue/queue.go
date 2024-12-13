package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cvhariharan/autopilot/internal/flow"
	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/lib/pq"
)

type Queue struct {
	store    repo.Store
	listener *pq.Listener
}

func NewQueue(r repo.Store) *Queue {
	return &Queue{
		store: r,
	}
}

func (q *Queue) Enqueue(ctx context.Context, f flow.Flow, input map[string]interface{}) (int32, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return -1, fmt.Errorf("error marshaling input to json: %w", err)
	}
	e, err := q.store.AddToQueue(ctx, repo.AddToQueueParams{
		FlowID: f.Meta.DBID,
		Input:  inputBytes,
	})

	if err != nil {
		return -1, fmt.Errorf("error adding flow to queue: %w", err)
	}

	return e.ID, err
}

func (q *Queue) ListenForExecution(ctx context.Context, count int) (<-chan int32, error) {
	if count <= 0 {
		count = 1
	}
	ch := make(chan int32, count)

	err := q.listener.Listen("new_flow")
	if err != nil {
		return nil, fmt.Errorf("error starting postgres listener for notifications: %w", err)
	}

	go func() error {
		select {
		case <-q.listener.Notify:
			// Dequeue the flow from the queue and execute it
		}
		return nil
	}()
	return ch, nil
}
