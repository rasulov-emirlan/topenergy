package domains

import (
	"fmt"
	"reflect"

	"github.com/rasulov-emirlan/topenergy-interview/internal/domains/tasks"
	"github.com/rasulov-emirlan/topenergy-interview/pkg/logging"
)

type CommonDependencies struct {
	Log *logging.Logger
}

func (deps CommonDependencies) Validate() error {
	if deps.Log == nil {
		return DependencyError{
			Dependency:       "CommonDependencies.Log",
			BrokenConstraint: "can't be nil",
		}
	}

	return nil
}

type TasksDependencies struct {
	Repo tasks.Repository
}

func (deps TasksDependencies) Validate() error {
	if isNil(deps.Repo) {
		return DependencyError{
			Dependency:       "TasksDependencies.Repo",
			BrokenConstraint: "can't be nil",
		}
	}

	return nil
}

type DependencyError struct {
	Dependency       string
	BrokenConstraint string
}

func (e DependencyError) Error() string {
	return fmt.Sprintf("dependency: %s, broke constraint: %s", e.Dependency, e.BrokenConstraint)
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
