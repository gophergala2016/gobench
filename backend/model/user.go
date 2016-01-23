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
}

type User struct {
	db   *mgo.Database
	coll *mgo.Collection
}

func NewUser(db *mgo.Database) (*User, error) {
	u := &User{db: db, coll: db.C("User")}
	return u, nil
}

func (u *User) CreateUser(ur *UserRow) error {
	_, err := u.coll.Upsert(bson.M{"login": ur.Login}, ur)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetByLogin(login string) (UserRow, error) {
	ur := new(UserRow)
	if err := u.coll.Find(bson.M{"login": login}).One(&ur); err != mgo.ErrNotFound {
		return UserRow{}, err
	}
	return *ur, nil
}
