package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// UserRow holds user's attributes
type UserRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// User's login
	Login string

	// User's token
	Token string

	// Link to avatar
	AvatarURL string

	// Packages followed by user
	Packages []bson.ObjectId
}

// User provides single point of access to user profile data
type User struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewUser creates model User
func NewUser(db *mgo.Database) (*User, error) {
	u := &User{db: db, coll: db.C("user")}
	return u, nil
}

// CreateOrUpdate creates if not exist or update user in db
func (u *User) CreateOrUpdate(ur *UserRow) (*mgo.ChangeInfo, error) {
	mu, err := u.GetByLogin(ur.Login)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}
	if err != mgo.ErrNotFound {
		mu.Token = ur.Token
		mu.AvatarURL = ur.AvatarURL
	} else {
		mu = ur
	}
	return u.Upsert(mu)
}

// Upsert upsert User in db
func (u *User) Upsert(ur *UserRow) (*mgo.ChangeInfo, error) {
	ci, err := u.coll.Upsert(bson.M{"login": ur.Login}, ur)
	if err != nil {
		return &mgo.ChangeInfo{}, err
	}
	return ci, nil
}

// GetByLogin find user in db by login
func (u *User) GetByLogin(login string) (*UserRow, error) {
	ur := new(UserRow)
	if err := u.coll.Find(bson.M{"login": login}).One(&ur); err != nil {
		return &UserRow{}, err
	}
	return ur, nil
}
