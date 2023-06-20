package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/rasulov-emirlan/topenergy-interview/config"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/health"
)

const servicePrefix = "tasks"

type RepoCombiner struct {
	tasks TasksRepo
}

func NewRepoCombiner(ctx context.Context, cfg config.Config) (RepoCombiner, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
	})

	res := rdb.Ping(ctx)
	if res.Err() != nil {
		return RepoCombiner{}, res.Err()
	}

	return RepoCombiner{
		tasks: TasksRepo{
			rdb: rdb,
		},
	}, nil
}

func (r RepoCombiner) Tasks() TasksRepo {
	return r.tasks
}

func (r RepoCombiner) Close() error {
	return r.tasks.rdb.Close()
}

func (r RepoCombiner) Check(ctx context.Context) health.Check {
	res := health.Check{
		Name:     "redis",
		Status:   health.StatusUP,
		Critical: true,
	}

	if err := r.tasks.rdb.Ping(ctx).Err(); err != nil {
		res.Status = health.StatusDOWN
		res.Message = err.Error()
	}

	return res
}
