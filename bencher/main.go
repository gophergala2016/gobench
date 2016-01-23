package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	AuthKey string `json: "authKey"`
	Email   string `json: "email"`
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	// reads config from file
	config, err := ReadConfig("./config.json")
	if err != nil {
		fmt.Println("Config file reading error. Details: ", err)
		return 1
	}

	// creates Logger (can be replaced with other popular logger)
	log := log.New(os.Stdout, "BEN ", log.Ldate|log.Ltime)
	log.Println("Application launched")

	br, err := NewBenchRunner(config.AuthKey, config.Email, log)
	if err != nil {
		log.Println("BenchRunner init failed. Details: ", err)
		return 2
	}

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
