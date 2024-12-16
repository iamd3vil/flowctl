package core

import (
	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Core struct {
	redisClient redis.UniversalClient
	store       repo.Store
	q           *asynq.Client
	flows       map[string]models.Flow
}

func NewCore(flows map[string]models.Flow, s repo.Store, q *asynq.Client, redisClient redis.UniversalClient) *Core {
	return &Core{
		store:       s,
		redisClient: redisClient,
		q:           q,
		flows:       flows,
	}
}
