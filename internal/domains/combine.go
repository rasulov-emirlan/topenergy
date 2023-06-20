package domains

import "github.com/rasulov-emirlan/topenergy-interview/internal/domains/tasks"

type DomainCombiner struct {
	tasksService tasks.Service
}

func NewDomainCombiner(commonDep CommonDependencies, tasksDep TasksDependencies) (DomainCombiner, error) {
	if err := commonDep.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	if err := tasksDep.Validate(); err != nil {
		return DomainCombiner{}, err
	}

	t := tasks.NewService(tasksDep.Repo, commonDep.Log)

	return DomainCombiner{
		tasksService: t,
	}, nil
}

func (c DomainCombiner) TasksService() tasks.Service {
	return c.tasksService
}
