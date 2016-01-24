package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)
//Config configs for backend
type Config struct {
	AuthKey 	string `json: "authKey"`
	Email   	string `json: "email"`
	BaseUrl		string `json: "baseurl"`
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	// important to have GOPATH defined
	if _, ok := os.LookupEnv("GOPATH"); !ok {
		fmt.Println("Environment variabe GOPATH not found!")
		return 1
	}

	// reads config from file
	config, err := ReadConfig("./config.json")
	if err != nil {
		fmt.Println("Config file reading error. Details: ", err)
		return 2
	}

	// creates Logger (can be replaced with other popular logger)
	log := log.New(os.Stdout, "BEN ", log.Ldate|log.Ltime)
	log.Println("Application launched")

	br, err := NewBenchClient(config.AuthKey, config.Email, config.BaseUrl, log)
	if err != nil {
		log.Println("BenchRunner init failed. Details: ", err)
		return 3
	}

	// Ctrl+C grcefully terminates application
	br.Run()

	log.Println("Application succesfully terminated")
	return 0
}

// ReadConfig reads application configuration from JSON file
func ReadConfig(fname string) (*Config, error) {

	var c Config
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
