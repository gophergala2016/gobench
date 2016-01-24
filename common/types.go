package common

import (
	"golang.org/x/tools/benchmark/parse"
)

type TaskRequest struct {
	AuthKey string `json: "authKey"`
	Email   string `json: "email"`
}

type TaskResponse struct {
	Id         string `json:"id"`
	PackageUrl string `json:"packageUrl"`

	// Type specifies task type.
	// TODO: support different task types: Benchmark, Build, Vet, etc.
	Type string `json:"type"`
}

type TaskResult struct {
	TaskRequest
	Id            string `json:"id"`
	Specification string `json: "specification"`

	// Round holds parsed bencmark results per GoMaxProcs 1-8
	Round      map[string]parse.Set
	BuildError string `json:"buildError"`
}