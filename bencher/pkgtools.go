package main

import (
	"errors"
	"os"
	"os/exec"
)

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

	cmd := exec.Command("go", "test", "-i", name)
	err := cmd.Start()
	if err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

/* Uncomment and on production platform
func cleanPackages() error {

	if debug {
		return nil
	}


		gopath, ok := os.LookupEnv("GOPATH")
		if !ok {
			return ErrGoPathNotFound
		}

		err := os.RemoveAll(gopath + "/src/")
		if err != nil {
			return err
		}

		return os.Mkdir(gopath+"/src/", os.ModeDir)
	return nil
}
*/
