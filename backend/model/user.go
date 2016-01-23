package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type UserRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// User login
	Login string

	// User token
	Token string

	// User avatar
	AvatarURL string

	Repos []string
}

type User struct {
	db   *mgo.Database
	coll *mgo.Collection
}

func NewUser(db *mgo.Database) (*User, error) {
	u := &User{db: db, coll: db.C("User")}
	return u, nil
}

func (u *User) CreateUser(ur *UserRow) (*mgo.ChangeInfo, error) {
	ci, err := u.coll.Upsert(bson.M{"login": ur.Login}, ur)
	if err != nil {
		return &mgo.ChangeInfo{}, err
	}
	return ci, nil
}

func (u *User) GetByLogin(login string) (*UserRow, error) {
	ur := new(UserRow)
	if err := u.coll.Find(bson.M{"login": login}).One(&ur); err != nil {
		return &UserRow{}, err
	}
	return ur, nil
}
