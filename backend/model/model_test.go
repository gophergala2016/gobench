package model_test

import (
	"encoding/json"
	"github.com/gophergala2016/gobench/backend/model"
	"testing"
)

func TestBenchmarkResult_DummyItems(t *testing.T) {

	br := model.BenchmarkResult{}

	items, _ := br.DummyItems("")
	buf, _ := json.MarshalIndent(items, "", "   ")
	t.Log(string(buf))
}
