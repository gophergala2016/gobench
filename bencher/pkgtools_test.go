package main

import (
    "testing"

)

func TestDownloadPackage( t *testing.T ) {
    path, err := downloadPackage( "github.com/gocraft/web" )
    if err != nil {
	t.Error( err )
    }
    t.Log( "Path: ", path )
}

func TestDownloadPackageDependencies( t *testing.T ) {
    err := downloadPackageDependencies( "github.com/gocraft/web" )
    if err != nil {
	t.Error( err )
    }
}

