package main

import (
	"errors"
	"os"
	"os/exec"
)

// stores errors
var (
	ErrGoPathNotFound = errors.New("Environment variable GOPATH not found")
)

func downloadPackage(name string) (string, error) {

	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return "", ErrGoPathNotFound
	}

	cmd := exec.Command("go", "get", "-u", name)
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	if err = cmd.Wait(); err != nil {
		return "", err
	}

	return gopath + "/src/" + name, nil
}

func downloadPackageDependencies(name string) error {

	return nil
}

func cleanPackages() error {

	if debug {
		return nil
	}

	/*	Uncomment on production only
		gopath, ok := os.LookupEnv("GOPATH")
		if !ok {
			return ErrGoPathNotFound
		}

		err := os.RemoveAll(gopath + "/src/")
		if err != nil {
			return err
		}

		return os.Mkdir(gopath+"/src/", os.ModeDir)
	*/
	return nil
}
