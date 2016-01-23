package model

import (
	"fmt"
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

	// PackageCount holds amount of packages in the repository
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

func (r *Repository) Add(repo RepositoryRow) error {
	_, err := r.coll.Upsert(bson.M{"url": repo.Url}, r)
	if err != nil {
		return err
	}
	return nil
}

// Items returns all repositories
func (r *Repository) Items(ids []string) ([]RepositoryRow, error) {
	oids := make([]bson.ObjectId, len(ids))
	for i := range ids {
		oids[i] = bson.ObjectIdHex(ids[i])
	}
	items := make([]RepositoryRow, 0)
	if err := r.coll.Find(bson.M{"_id": bson.M{"$in": oids}}).All(&items); err != nil {
		fmt.Println(items)
		return nil, err
	}
	return items, nil
}