package common

import (
	"golang.org/x/tools/benchmark/parse"
)

// TaskRequest stores info about task's author
type TaskRequest struct {
	AuthKey string `json: "authKey"`
	Email   string `json: "email"`
}

// TaskResponse stores responses of tasks
type TaskResponse struct {
	Id          string `json:"id"`
	PackageName string `json:"packageName"`

	// Type specifies task type.
	// TODO: support different task types: Benchmark, Build, Vet, etc.
	Type []string `json:"type"`
}

// TaskResult stores result(errors) of task
type TaskResult struct {
	TaskRequest
	Id            string `json:"id"`
	Specification string `json: "specification"`

	// Round holds parsed bencmark results per GoMaxProcs "cpu1", "cpu2", "cpu4"
	Round      map[string]parse.Set
	BuildError string `json:"buildError"`
}
