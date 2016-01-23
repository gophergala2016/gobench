package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// TestEnvironmentRow holds single test environment attributes
type TestEnvironmentRow struct {
	Id bson.ObjectId

	// AuthKey used by remote gobech client to retrive next task
	// assigned to it
	AuthKey string

	// Name presents on UI
	Name string

	// Specification holds description of testing environment
	// filled manually or after the 1st execution
	Specification string

	// LastSpecification holds description of testing environment arrived
	// together with last report
	LastSpecification string

	// Weight used for sorting. Small or slow test environments has less value
	Weight int
}

// TestEnvironment is single point of access
type TestEnvironment struct {
	db   *mgo.Database
	coll *mgo.Collection
}

func NewTestEnvironment(db *mgo.Database) (*TestEnvironment, error) {
	ws := &TestEnvironment{db: db, coll: db.C("TestEnv")}
	return ws, nil
}

// Items returns all test environments sorted by Weight (asc)
func (te *TestEnvironment) Items() ([]TestEnvironmentRow, error) {
	items := make([]TestEnvironmentRow, 0)
	err := te.coll.Find(nil).Sort("Weight").All(&items)
	return items, err
}
