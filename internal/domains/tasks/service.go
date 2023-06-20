package tasks

import (
	"context"
	"errors"

	"github.com/rasulov-emirlan/topenergy-interview/pkg/logging"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/rasulov-emirlan/topenergy-interview/internal/domains/tasks"

type (
	Repository interface {
		Create(ctx context.Context, task Task) (Task, error)
		Read(ctx context.Context, id string) (Task, error)
		ReadAll(ctx context.Context) ([]Task, error)
		Update(ctx context.Context, task Task) (Task, error)
		Delete(ctx context.Context, id string) error
	}

	Service interface {
		Create(ctx context.Context, task Task) (Task, error)
		Read(ctx context.Context, id string) (Task, error)
		ReadAll(ctx context.Context) ([]Task, error)
		Update(ctx context.Context, task Task) (Task, error)
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo Repository
		log  *logging.Logger
	}
)

var _ Service = (*service)(nil)

func NewService(repo Repository, log *logging.Logger) service {
	return service{
		repo: repo,
		log:  log,
	}
}

func (s service) Create(ctx context.Context, task Task) (Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "tasks.Create")
	defer span.End()
	defer s.log.Sync()

	t, err := s.repo.Create(ctx, task)
	if err != nil {
		s.log.Error("tasks.Create", logging.String("stage", "db"), logging.Error("err", err))
		return Task{}, errors.New("failed to create task")
	}
	s.log.Info("tasks.Create", logging.String("id", t.ID))
	return t, nil
}

func (s service) Read(ctx context.Context, id string) (Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "tasks.Read")
	defer span.End()
	defer s.log.Sync()

	t, err := s.repo.Read(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.log.Debug("tasks.Read", logging.String("stage", "db"), logging.Error("err", err))
			return Task{}, ErrTaskNotFound
		}
		s.log.Error("tasks.Read", logging.String("stage", "db"), logging.Error("err", err))
		return Task{}, errors.New("failed to read task")
	}
	s.log.Info("tasks.Read", logging.String("id", t.ID))
	return t, nil
}

func (s service) ReadAll(ctx context.Context) ([]Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "tasks.ReadAll")
	defer span.End()
	defer s.log.Sync()

	tasks, err := s.repo.ReadAll(ctx)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.log.Debug("tasks.ReadAll", logging.String("stage", "db"), logging.Error("err", err))
			return nil, ErrTaskNotFound
		}
		s.log.Error("tasks.ReadAll", logging.String("stage", "db"), logging.Error("err", err))
		return nil, errors.New("failed to read tasks")
	}
	s.log.Info("tasks.ReadAll", logging.Int("count", len(tasks)))
	return tasks, nil
}

func (s service) Update(ctx context.Context, task Task) (Task, error) {
	ctx, span := otel.Tracer(otelName).Start(ctx, "tasks.Update")
	defer span.End()
	defer s.log.Sync()

	t, err := s.repo.Update(ctx, task)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.log.Debug("tasks.Update", logging.String("stage", "db"), logging.Error("err", err))
			return Task{}, ErrTaskNotFound
		}
		s.log.Error("tasks.Update", logging.String("stage", "db"), logging.Error("err", err))
		return Task{}, errors.New("failed to update task")
	}
	s.log.Info("tasks.Update", logging.String("id", t.ID))
	return t, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Tracer(otelName).Start(ctx, "tasks.Delete")
	defer span.End()
	defer s.log.Sync()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.log.Debug("tasks.Delete", logging.String("stage", "db"), logging.Error("err", err))
			return ErrTaskNotFound
		}
		s.log.Error("tasks.Delete", logging.String("stage", "db"), logging.Error("err", err))
		return errors.New("failed to delete task")
	}
	s.log.Info("tasks.Delete", logging.String("id", id))
	return nil
}
