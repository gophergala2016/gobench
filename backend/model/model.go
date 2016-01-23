package model

import (
	"labix.org/v2/mgo"
	"log"
)

// Model provides single point of access to all models
type Model struct {
	Repository *Repository
	Package    *Package
	logger     *log.Logger
}

// New creates object Model
func New(db *mgo.Database, l *log.Logger) (*Model, error) {

	var err error

	m := &Model{logger: l}

	m.Repository, err = NewRepository(db)
	if err != nil {
		return nil, err
	}

	m.Package, err = NewPackage(db)
	if err != nil {
		return nil, err
	}
	return m, nil
}
