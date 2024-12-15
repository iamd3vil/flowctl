package core

import (
	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Core struct {
	redisClient redis.UniversalClient
	store       repo.Store
	q           *asynq.Client
}

func NewCore(s repo.Store, q *asynq.Client, redisClient redis.UniversalClient) *Core {
	return &Core{
		store:       s,
		redisClient: redisClient,
		q:           q,
	}
}
