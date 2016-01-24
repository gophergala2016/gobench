package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
	"math/rand"
)

type RepositoryEngine string

const (
	Git RepositoryEngine = "git"

// TODO: Implement support after GopherGala
	Bazaar    RepositoryEngine = "bazaar"
	Mercurial RepositoryEngine = "mercurial"
)

// PackageRow holds package attributes
type PackageRow struct {
	Id            bson.ObjectId `bson:"_id,omitempty"`

	// Name of package in a gopher way (labix.org/v2/mgo)
	Name          string `bson:"name"`

	// Url holds full package url
	Url           string `bson: "url"`

	// Description of the package
	Description   string `bson:"description"`

	// All tags of This repository
	Tags []RepositoryTag `bson:"tags"`

	// Repository holds repository url (https://github.com or https://labix.org, etc)
	RepositoryUrl string `bson:"repositoryUrl"`

	// Repository's engine
	Engine        RepositoryEngine `bson:"engine"`

	// Created holds time
	Created       time.Time

	// Created holds time of the last update
	Updated       time.Time

	// LastCommitUid holds hash of the the last commit
	LastCommitId  string
}

type RepositoryTag struct {
	Name string `bson: "name"`
	Zip string `bson: "zip"`
	Tar string `bson: "tar"`
	Commit string `bson: "commit"`
}

// Package provides single point of access to all packages
type Package struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewPackage creates package model
func NewPackage(db *mgo.Database) (*Package, error) {
	p := &Package{db: db, coll: db.C("package")}
	idx, err := p.coll.Indexes()
	if err != nil {
		return nil, err
	}
	for i := range idx {
		if len(idx[i].Key) > 0 && idx[i].Key[0] == "url" {
			return p, nil
		}
	}

	return p, p.coll.EnsureIndex(mgo.Index{Key: []string{"url"}, Unique: true, DropDups: true})
}

// Add inserts new package and ignores if package exist already
func (p *Package) Add(pr *PackageRow) error {
	pr.Created = time.Now()

	err := p.coll.Insert(pr)
	if err != nil && !mgo.IsDup(err) {
		return err
	}
	return nil
}

func (p *Package) GetItem(name string) (PackageRow, error) {
	item := PackageRow{}
	if err := p.coll.Find(bson.M{"url": bson.RegEx{name, ""}}).One(&item); err != nil {
		return item, err
	}
	return item, nil
}

func (p *Package) GetItems(name string) ([]PackageRow, error) {
	item := make([]PackageRow, 0)
	if err := p.coll.Find(bson.M{"url": bson.RegEx{name, ""}}).All(&item); err != nil {
		return nil, err
	}
	return item, nil
}

// All returns all packages
func (p *Package) All() ([]PackageRow, error) {
	items := make([]PackageRow, 0)
	if err := p.coll.Find(nil).All(&items); err != mgo.ErrNotFound {
		return nil, err
	}
	return items, nil
}

// Favorites returns packages starred by github user
func (p *Package) Favorites(userName string) ([]PackageRow, error) {
	items := make([]PackageRow, 0)
	// TODO: добавить условие поиска/фильтрации
	if err := p.coll.Find(nil).All(&items); err != mgo.ErrNotFound {
		return nil, err
	}
	return items, nil
}

// Items returns all repositories
func (p *Package) GetItemsByIdSlice(oids []bson.ObjectId) ([]PackageRow, error) {
	items := make([]PackageRow, 0)
	if err := p.coll.Find(bson.M{"_id": bson.M{"$in": oids}}).All(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Package) Items(oids []bson.ObjectId) ([]PackageRow, error) {
	return p.GetItemsByIdSlice(oids)
}


func (p *Package) DummyList() ([]PackageRow, error) {
	rInt := rand.Intn(10)
	items := make([]PackageRow, 0)
	if err := p.coll.Find(bson.M{}).Skip(rInt).Limit(10).All(&items); err != nil {
		return nil, err
	}
	return items, nil
}
