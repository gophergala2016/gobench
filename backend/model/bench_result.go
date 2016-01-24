package model

import (
	"golang.org/x/tools/benchmark/parse"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// BenchmarkResultRow holds package benchmark
type BenchmarkResultRow struct {
	Id bson.ObjectId `bson:"_id,omitempty"`

	// PackageUrl holds full URL to package (github.com/gorilla/session)
	PackageName string `bson:"packageName"`

	// Created holds row saving data. Used for chart sorting
	Created time.Time `bson:"created"`

	// Value holds benchmark values
	Value map[string]parse.Set `bson:"value"`

	// TestEnvSpecification holds description of HW/OS where task executed
	TestEnvSpecification string `bson:"testEnv"`
}

// BenchmarkResult provides single point of access to all benchmarking results
type BenchmarkResult struct {
	db   *mgo.Database
	coll *mgo.Collection
}

// NewBenchmarkResult creates BenchmarkResult model
func NewBenchmarkResult(db *mgo.Database) (*BenchmarkResult, error) {
	t := &BenchmarkResult{db: db, coll: db.C("benchmarkResult")}
	return t, nil
}

// Add saves bechmark results
func (br *BenchmarkResult) Add(pkgName, testEnv string, value map[string]parse.Set) error {
	item := BenchmarkResultRow{Created: time.Now(), PackageName: pkgName, TestEnvSpecification: testEnv, Value: value}
	if err := br.coll.Insert(item); err != nil {
		return err
	}
	return nil
}

// Items retrives benchmark results for specific package identified by url
func (br *BenchmarkResult) Items(pkgName string) ([]BenchmarkResultRow, error) {
	var items []BenchmarkResultRow

	if err := br.coll.Find(bson.M{"packageName": pkgName}).Sort("created").All(&items); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return items, nil
}

// DummyItems returns bechmark results Dummy Items
func (br *BenchmarkResult) DummyItems(url string) ([]BenchmarkResultRow, error) {
	items := make([]BenchmarkResultRow, 3)

	items[0].Created = time.Now().Add(-1 * time.Hour * 24)
	items[1].Created = time.Now().Add(-1 * time.Hour * 12)
	items[2].Created = time.Now()

	items[0].Value = make(map[string]parse.Set)
	items[1].Value = make(map[string]parse.Set)
	items[2].Value = make(map[string]parse.Set)

	items[0].Value["cpu1"] = make(parse.Set)
	items[0].Value["cpu2"] = make(parse.Set)
	items[0].Value["cpu4"] = make(parse.Set)

	items[1].Value["cpu1"] = make(parse.Set)
	items[1].Value["cpu2"] = make(parse.Set)
	items[1].Value["cpu4"] = make(parse.Set)

	items[2].Value["cpu1"] = make(parse.Set)
	items[2].Value["cpu2"] = make(parse.Set)
	items[2].Value["cpu4"] = make(parse.Set)

	items[0].Value["cpu1"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 10, AllocedBytesPerOp: 200}}
	items[0].Value["cpu2"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 20, AllocedBytesPerOp: 300}}
	items[0].Value["cpu4"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 30, AllocedBytesPerOp: 400}}

	items[0].Value["cpu1"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 10, AllocedBytesPerOp: 200}}
	items[0].Value["cpu2"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 20, AllocedBytesPerOp: 300}}
	items[0].Value["cpu4"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 30, AllocedBytesPerOp: 400}}

	items[1].Value["cpu1"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 10, AllocedBytesPerOp: 200}}
	items[1].Value["cpu2"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 20, AllocedBytesPerOp: 300}}
	items[1].Value["cpu4"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 30, AllocedBytesPerOp: 400}}

	items[1].Value["cpu1"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 10, AllocedBytesPerOp: 200}}
	items[1].Value["cpu2"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 20, AllocedBytesPerOp: 300}}
	items[1].Value["cpu4"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 30, AllocedBytesPerOp: 400}}

	items[2].Value["cpu1"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 50, AllocedBytesPerOp: 200}}
	items[2].Value["cpu2"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 70, AllocedBytesPerOp: 1300}}
	items[2].Value["cpu4"]["BenchmarkInsertRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 100, AllocedBytesPerOp: 2400}}

	items[2].Value["cpu1"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 10, AllocedBytesPerOp: 200}}
	items[2].Value["cpu2"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 20, AllocedBytesPerOp: 300}}
	items[2].Value["cpu4"]["BenchmarkUpdateRow"] = []*parse.Benchmark{&parse.Benchmark{NsPerOp: 30, AllocedBytesPerOp: 400}}

	return items, nil
}
