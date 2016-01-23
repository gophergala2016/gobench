package main

import (
	"fmt"
	"github.com/gophergala2016/gobench/backend"
	"github.com/gophergala2016/gobench/common/config"
	"github.com/gophergala2016/gobench/frontend"
	"log"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {

	// reads and parses application config params
	config, err := config.ReadFromFile("./config.json")
	if err != nil {
		fmt.Println("Config file reading error. Details: ", err)
		return 1
	}

	// creates Logger (can be replaced with other popular logger)
	log := log.New(os.Stdout, "GOB ", log.Ldate|log.Ltime)
	log.Println("Application launched")

	// creates Backend object
	back, err := backend.New(&config.Backend, log)
	if err != nil {
		log.Println("Backend initialisation error. Details: ", err)
		return 2
	}

	// creates Frontend objects
	front, err := frontend.New(&config.Frontend, log, back)
	if err != nil {
		log.Println("Frontend initialisation error. Details: ", err)
		return 3
	}

	// start background processes
	if err := back.Start(); err != nil {
		log.Println("Backend processes launching error. Details: ", err)
		return 4
	}

	// start HTTP listener and handlers
	if err := front.Start(); err != nil {
		log.Println("Fronend listeners/processes error. Details: ", err)
		return 5
	}

	log.Println("Application succesfully terminated")
	return 0
}
