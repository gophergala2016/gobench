package model_test

import (
	"encoding/json"
	"github.com/gophergala2016/gobench/backend"
	"github.com/gophergala2016/gobench/backend/model"
	"log"
	"os"
	"testing"
	//	"time"
)

var cfg = backend.Config{
	Mongo: backend.DatabaseConfig{IP: "127.0.0.1", Port: 27017, Name: "gobench"},
	Debug: true,
}

func aTestBenchmarkResult_DummyItems(t *testing.T) {

	br := model.BenchmarkResult{}

	items, _ := br.DummyItems("")
	buf, _ := json.MarshalIndent(items, "", "   ")
	t.Log(string(buf))

}

func TestBenchmarkResult_DummyItems(t *testing.T) {

	back, err := backend.New(&cfg, log.New(os.Stdout, "TEST ", 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(back.Model.TestEnvironment.Items())

	pkgs, err := back.Model.Package.All()
	if err != nil {
		t.Fatal(err)
	}

	/*
		err = back.Model.Package.Add(&model.PackageRow{
			Name:          "github.com/bradfitz/slice",
			Url:           "https://github.com/bradfitz/slice",
			RepositoryUrl: "https://github.com",
			Engine:        "git",
			Created:       time.Now(),
			Updated:       time.Now()})

	*/
	pkg, err := back.Model.Package.GetItem("github.com/bradfitz/slice")
	if pkg.Name != "github.com/bradfitz/slice" {
		t.Fatal(err)
	}
	t.Log(pkg, err)

	task, err := back.Model.Task.Next("change-secret-1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(task)
}
