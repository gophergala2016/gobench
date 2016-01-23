package model

import (
	"labix.org/v2/mgo"
	"log"
)

// Model provides single point of access to all models
type Model struct {
	logger *log.Logger
}

// New creates object Model
func New(db *mgo.Database, l *log.Logger) (*Model, error) {

	m := &Model{logger: l}

	return m, nil
}
