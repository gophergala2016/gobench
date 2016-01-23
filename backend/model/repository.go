package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type RepositoryEngine string

const (
	Git RepositoryEngine = "git"

	// TODO: Implement support after GopherGala
	Bazaar    RepositoryEngine = "bazaar"
	Mercurial RepositoryEngine = "mercurial"
)

// Repository holds repository attributes
type RepositoryRow struct {
	Id bson.ObjectId

	// Name of repository
	Name string

	// Url of repository
	Url string

	// Repository's engine
	Engine RepositoryEngine

	// PackageCount holds amount of packages stored in the repository
	PackageCount int `bson:"-"`
}

// Repository provides single point of access to all repositories
type Repository struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewRepository creates repository model
func NewRepository(db *mgo.Database) (*Repository, error) {
	r := &Repository{db: db, coll: db.C("Repository")}
	return r, nil
}

// Items returns all repositories
func (r *Repository) Items() ([]RepositoryRow, error) {
	items := make([]RepositoryRow, 0)
	if err := r.coll.Find(nil).All(&items); err != mgo.ErrNotFound {
		return nil, err
	}
	return items, nil
}
