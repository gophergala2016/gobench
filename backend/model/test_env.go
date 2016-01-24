package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// TestEnvironmentRow holds single test environment attributes
type TestEnvironmentRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// AuthKey used by remote gobech client to retrive next task
	// assigned to it
	AuthKey string `bson:"authKey"`

	// Name presents on UI
	Name string `bson:"name"`

	// Specification holds description of testing environment
	// filled manually or after the 1st execution
	Specification string `bson:"specification"`

	// LastSpecification holds description of testing environment arrived
	// together with last report
	LastSpecification string

	// Weight used for sorting. Small or slow test environments has less value
	Weight int `bson:"weight"`
}

// TestEnvironment is single point of access
type TestEnvironment struct {
	db   *mgo.Database
	coll *mgo.Collection
}
// NewTestEnvironment create and return new environment
func NewTestEnvironment(db *mgo.Database) (*TestEnvironment, error) {
	ws := &TestEnvironment{db: db, coll: db.C("testEnv")}
	return ws, nil
}

// Items returns all test environments sorted by Weight (asc)
func (te *TestEnvironment) Items() ([]TestEnvironmentRow, error) {
	var items []TestEnvironmentRow
	err := te.coll.Find(nil).Sort("weight").All(&items)
	return items, err
}

// Exist checks existense of test environment by authKey
func (te *TestEnvironment) Exist(authKey string) (bool, error) {
	cnt, err := te.coll.Find(bson.M{"authKey": authKey}).Count()
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}
	return cnt > 0, nil
}
