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
}

type TaskResult struct {
	TaskRequest
	Id            string `json:"id"`
	Specification string `json: "specification"`

	// Result holds parsed bencmark results per GoMaxProc 1-8
	Result     map[string]parse.Set
	BuildError string `json:"buildError"`
}
