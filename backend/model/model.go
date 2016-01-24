package model

import (
	"errors"
	"labix.org/v2/mgo"
	"log"
)

var (
	ErrAuthKey  = errors.New("Wrong authKey value!")
	ErrNotFound = errors.New("Not found!")
)

// Model provides single point of access to all models
type Model struct {
	Package         *Package
	User            *User
	TestEnvironment *TestEnvironment
	Task            *Task
	BenchmarkResult *BenchmarkResult
	logger          *log.Logger
}

// New creates object Model
func New(db *mgo.Database, l *log.Logger) (*Model, error) {

	var err error

	m := &Model{logger: l}

	m.Package, err = NewPackage(db)
	if err != nil {
		return nil, err
	}

	m.User, err = NewUser(db)
	if err != nil {
		return nil, err
	}

	m.TestEnvironment, err = NewTestEnvironment(db)
	if err != nil {
		return nil, err
	}

	m.Task, err = NewTask(db)
	if err != nil {
		return nil, err
	}

	m.BenchmarkResult, err = NewBenchmarkResult(db)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// PackageName registers individual task for each test environment
func (m *Model) RegisterTasks(pkgName string) error {

	te, err := m.TestEnvironment.Items()
	if err != nil {
		return err
	}

	for i := range te {
		err := m.Task.Register(pkgName, te[i].AuthKey, []string{"becnmark"})
		if err != nil {
			m.logger.Printf("Could not register new task for package %s. Details: %s", pkgName, err)
		}
	}

	return nil
}
