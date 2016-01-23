package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// TaskRow holds task description
type TaskRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// PackageUrl holds full URL to package (github.com/gorilla/session)
	PackageUrl string `bson:"packageUrl"`

	// AuthKey holds key identifies testing environment
	AuthKey string `bson:"authKey"`

	// Created holds row creation time. Uses for sorting
	Created time.Time
}

// Task provides single point of access to all bench tasks to execute
type Task struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewTask creates Task model
func NewTask(db *mgo.Database) (*Task, error) {
	t := &Task{db: db, coll: db.C("Task")}
	return t, nil
}

// Next returns next task for test environment identified by authKey
func (t *Task) Next(authKey string) (*TaskRow, error) {

	var tr TaskRow
	err := t.coll.Find(bson.M{"authKey": authKey}).One(&tr) // .Sort("created").
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &tr, nil
}
