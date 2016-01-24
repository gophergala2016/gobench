package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// UserRow stores User info
type UserRow struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`

	// User login
	Login     string

	// User token
	Token     string

	// User avatar
	AvatarURL string

	Packages  []bson.ObjectId
}
// User stores mongo's info about user collection
type User struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewUser add new user in db collection
func NewUser(db *mgo.Database) (*User, error) {
	u := &User{db: db, coll: db.C("user")}
	return u, nil
}

// CreateOrUpdateUser create or, if it exists, update user in db
func (u *User) CreateOrUpdateUser(ur *UserRow) (*mgo.ChangeInfo, error) {
	mu, err := u.GetByLogin(ur.Login)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}
	if (err != mgo.ErrNotFound) {
		mu.Token = ur.Token
		mu.AvatarURL = ur.AvatarURL
	} else {
		mu = ur
	}
	return u.UpsertUser(mu)
}

// UpsertUser upsert User in db
func (u *User) UpsertUser(ur *UserRow) (*mgo.ChangeInfo, error) {
	ci, err := u.coll.Upsert(bson.M{"login": ur.Login}, ur)
	if err != nil {
		return &mgo.ChangeInfo{}, err
	}
	return ci, nil
}

//GetByLogin find user in db by login
func (u *User) GetByLogin(login string) (*UserRow, error) {
	ur := new(UserRow)
	if err := u.coll.Find(bson.M{"login": login}).One(&ur); err != nil {
		return &UserRow{}, err
	}
	return ur, nil
}
