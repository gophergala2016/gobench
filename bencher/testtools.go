package main

import (
//    "errors"
//    "os"
    "os/exec"
//    "encoding/json"
    "golang.org/x/tools/benchmark/parse"
    "log"
)

func runTest( name string, maxprocs int ) ( parse.Set, error ) {
    log.Printf( "Start testing of %s", name )
    log.Print( "Create cmd............." )
    cmd := exec.Command( "go", "test", "-bench=.", "-benchmem", "-run=NONE", name + "/..." )
    log.Println ( "Ok" )

    log.Print( "Create StdoutPipe......" )
    stdout, err := cmd.StdoutPipe()
    if err != nil {
	log.Println( "Get pipe failed:", err )
	return nil, err
    }
    log.Println ( "Ok" )

    log.Print( "Start cmd.............." )
    if err := cmd.Start(); err != nil {
	log.Println( "Execution of ", name, " failed:", err )
	return nil, err
    }
    
    log.Println ( "Ok" )


    log.Print( "Start ParseSet........" )
    res, err := parse.ParseSet( stdout )
    
    if err != nil {
	log.Println( "Parsing of output failed ", name, " failed:", err )
	return nil, err
    }
    log.Println ( "Ok" )

    log.Print( "Wait for exit........" )
    if err := cmd.Wait(); err != nil {
	log.Print ("Failed:", err )
	return nil, err
    }

    log.Println ( "Ok" )

    return res, nil
    
}

