package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rasulov-emirlan/topenergy-interview/internal/domains/tasks"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/rasulov-emirlan/topenergy-interview/internal/storage/redis"

type TasksRepo struct {
	rdb *redis.Client
}

func (r TasksRepo) Create(ctx context.Context, task tasks.Task) (tasks.Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "TasksRepo.Create")
	defer span.End()

	task.ID = uuid.New().String()
	key := fmt.Sprintf("%s:%s", servicePrefix, task.ID)
	return task, r.rdb.HSet(ctx, key, "title", task.Title, "description", task.Description).Err()
}

func (r TasksRepo) Read(ctx context.Context, id string) (tasks.Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "TasksRepo.Read")
	defer span.End()

	key := fmt.Sprintf("%s:%s", servicePrefix, id)
	res, err := r.rdb.HGetAll(ctx, key).Result()
	// return err not found
	if err != nil {
		if err == redis.Nil {
			return tasks.Task{}, fmt.Errorf("%w: %s", tasks.ErrTaskNotFound, id)
		}
		return tasks.Task{}, err
	}

	if len(res) == 0 {
		return tasks.Task{}, fmt.Errorf("%w: %s", tasks.ErrTaskNotFound, id)
	}

	return tasks.Task{
		ID:          id,
		Title:       res["title"],
		Description: res["description"],
	}, nil
}

func (r TasksRepo) ReadAll(ctx context.Context) ([]tasks.Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "TasksRepo.ReadAll")
	defer span.End()

	keys, err := r.rdb.Keys(ctx, fmt.Sprintf("%s:*", servicePrefix)).Result()
	if err != nil {
		return nil, err
	}
	var result []tasks.Task
	for _, key := range keys {
		res, err := r.rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		result = append(result, tasks.Task{
			ID:          strings.ReplaceAll(key, fmt.Sprintf("%s:", servicePrefix), ""), // TODO: refactor
			Title:       res["title"],
			Description: res["description"],
		})
	}
	return result, nil
}

func (r TasksRepo) Update(ctx context.Context, task tasks.Task) (tasks.Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "TasksRepo.Update")
	defer span.End()

	key := fmt.Sprintf("%s:%s", servicePrefix, task.ID)
	err := r.rdb.HSet(ctx, key, "title", task.Title, "description", task.Description).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return tasks.Task{}, fmt.Errorf("%w: %s", tasks.ErrTaskNotFound, task.ID)
		}
		return tasks.Task{}, err
	}
	return task, nil
}

func (r TasksRepo) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Tracer(otelName).Start(ctx, "TasksRepo.Delete")
	defer span.End()

	key := fmt.Sprintf("%s:%s", servicePrefix, id)
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%w: %s", tasks.ErrTaskNotFound, id)
		}
		return err
	}
	return nil
}
