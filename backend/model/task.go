package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// TaskRow holds task description
type TaskRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// PackageName holds full URL to package (github.com/gorilla/session)
	PackageName string `bson:"packageName"`

	// AuthKey holds key identifies testing environment
	AuthKey string `bson:"authKey"`

	// Creates holds row creation time
	Created time.Time `bson:"created"`

	// Type of task
	Type []string `bson:"type"`

	// Assigned holds tasks assignment time. Used for excluding
	Assigned time.Time `bson:"assigned"`
}

// Task provides single point of access to tasks
type Task struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewTask creates Task model, mongo's collection "task"
func NewTask(db *mgo.Database) (*Task, error) {
	t := &Task{db: db, coll: db.C("task")}
	return t, nil
}

// Next returns next task for test environment identified by authKey
func (t *Task) Next(authKey string) (*TaskRow, error) {

	// TODO: low priority, add sorting by Created
	var tr TaskRow
	err := t.coll.Find(bson.M{"authKey": authKey}).One(&tr)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	// TODO: medium priority, mark the task as taken and release in N minutes
	// if row still exists

	return &tr, nil
}

// Register creates task for each test environment
func (t *Task) Register(pkgName, authKey string, typ []string) error {
	tr := TaskRow{PackageName: pkgName, AuthKey: authKey, Type: typ, Created: time.Now()}
	return t.coll.Insert(&tr)
}

// Get retriives task by id
func (t *Task) Get(id string) (*TaskRow, error) {

	var tr TaskRow
	err := t.coll.FindId(bson.ObjectId(id)).One(&tr) // .Sort("created").
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &tr, nil
}

func (t *Task) Exist(id string) (bool, error) {
	_, err := t.Get(id)
	return err == nil, err
}

// Delete removes task from the list
func (t *Task) Delete(id string) error {
	err := t.coll.RemoveId(bson.ObjectIdHex(id))
	return err
}
