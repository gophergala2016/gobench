package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// PackageRow holds package attributes
type PackageRow struct {
	Id bson.ObjectId

	// Url
	Url string `bson: "url"`

	// Description of the package
	Description string `bson:"description"`

	// Repository's engine
	Engine RepositoryEngine `bson:"engine`

	// PackageCount holds amount of packages store in the repository
	PackageCount int `bson:"packagecount"`
}

// Package provides single point of access to all packages
type Package struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewPackage creates package model
func NewPackage(db *mgo.Database) (*Package, error) {
	r := &Package{db: db, coll: db.C("Package")}
	return r, nil
}

func (r *Package) GetItem(name string) ([]PackageRow, error) {
      item := make([]PackageRow,0)
      if err := r.coll.Find(bson.M{"url": bson.RegEx{name,""}}).All(&item); err != nil {
		  return nil, err
	  }
	 return item, nil
}

// Items returns all packages
func (r *Package) Items() ([]PackageRow, error) {
	items := make([]PackageRow, 0)
	if err := r.coll.Find(nil).All(&items); err != mgo.ErrNotFound {
		return nil, err
	}
	return items, nil
}

// Favorites returns packages starred by github user
func (r *Package) Favorites(userName string) ([]PackageRow, error) {
	items := make([]PackageRow, 0)
	// TODO: добавить условие поиска/фильтрации
	if err := r.coll.Find(nil).All(&items); err != mgo.ErrNotFound {
		return nil, err
	}
	return items, nil
}
