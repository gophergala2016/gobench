package job

import (
	"labix.org/v2/mgo"
	"log"
)

// Job provides single point of access to all background running jobs
type Job struct {
	logger *log.Logger
}

// New creates object Job
func New(db *mgo.Database, l *log.Logger) (*Job, error) {

	m := &Job{logger: l}

	return m, nil
}

func (j *Job) Start() error {

	return nil
}
