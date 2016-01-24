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
    bb, err := runTest ( "github.com/gocraft/web", 1)
//    bb, err := runTest ( "github.com/gin-gonic/gin", 1)
    if err != nil {
	t.Error( err )
    }
    t.Log ("---finished---")
    for key, value := range bb  {
	t.Log ( "Benchmark:", key );
	for _,value1 := range value {
	
	    t.Log ("{" );

	    t.Log( "  Name   : ", value1.Name )
	    t.Log( "  N      : ", value1.N )
	    t.Log( "  NsPerOp: ", value1.NsPerOp )
	    t.Log( "  AllocedBytesPerOp: ", value1.AllocedBytesPerOp )
	    t.Log( "  AllocsPerOp: ", value1.AllocsPerOp )
	    t.Log( "  MBPerS: ", value1.MBPerS )
	    t.Log( "  Measured: ", value1.Measured )
	    t.Log( "  Ord: ", value1.Ord )

	    t.Log ("}" );
	}
    }

}
