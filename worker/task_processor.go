package worker

import (
	"context"

	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	CriticalQ = "critical"
	DefaultQ  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

// Start implements TaskProcessor.
func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				CriticalQ: 10,
				DefaultQ:  5,
			},
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
