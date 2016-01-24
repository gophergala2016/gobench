package main

import (
//    "errors"
//    "os"
//    "os/exec"
//    "encoding/json"
//    "golang.org/x/tools/benchmark/parse"
//    "log"
    "testing"
)


func TestRunTest( t *testing.T ) {
//    bb, err := runTest ( "github.com/gocraft/web", 1)
    bb, err := runTest ( "github.com/gin-gonic/gin", 1)
    if err != nil {
	t.Error( err )
    }
    for key, value := range bb  {
	t.Log( "res[", key , "] = ", value )
    }

}
