package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gophergala2016/gobench/common"
	"golang.org/x/tools/benchmark/parse"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
	"strconv"
//	"fmt"
)

const (
//	baseUrl           = "http://127.0.0.1:8080"
	debug             = true
	nextTaskLocation       = "/api/task/next"
	submitResultLocation   = "/api/task/submit"
	maxSubmitAttepmts = 5
)

// BenchClient implements client interface to goben.ch server
type BenchClient struct {
	client        *http.Client
	log           *log.Logger
	authKey       string
	email         string
	baseUrl	      string
	nextTaskUrl   string
	submitResultUrl string
	stopCh        chan os.Signal
	specification string
}

// NewBenchClient creates BenchClient instance
func NewBenchClient(authKey, email, baseurl string, l *log.Logger) (*BenchClient, error) {

	br := &BenchClient{
		authKey: authKey,
		email:   email,
		baseUrl: baseurl,
		nextTaskUrl: baseurl + nextTaskLocation,
		submitResultUrl: baseurl + submitResultLocation,
		client:  &http.Client{Timeout: 2 * time.Second},
		log:     l,
	}

	// TODO: identify current machine specification: RAM, CPU, OS, etc.
	// save to br.specification
	br.specification = "Bare metal/Intel i5-5200K, 4 core, Ubuntu 14.04, RAM: 16G"

	err := br.Ping()
	if err != nil && debug == false {
		// exit with error in debug mode
		return nil, err
	}
	return br, nil
}

// Ping checks goben.ch availability
func (br *BenchClient) Ping() error {
	return br.ping()
}

func (br *BenchClient) ping() error {

	// TODO: process HTTP 301
	resp, err := br.client.Head(br.baseUrl)
	if err != nil {
		return err
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return err
}

// Run starts goben.ch client
func (br *BenchClient) Run() {

	br.stopCh = make(chan os.Signal)
	signal.Notify(br.stopCh, syscall.SIGINT)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		br.run()
	}()

	br.log.Println("Starting goben.ch client")
	select {
	case signal := <-br.stopCh:
		br.log.Println("Got signal: ", signal)
	}
	br.log.Println("Stopping client")
	br.stop()
	br.log.Println("Waiting task finalization")
	wg.Wait()
	return
}

func (br *BenchClient) run() {

	for {
		select {
		case <-br.stopCh:
			return
		default:
			//If the channel is still open, continue as normal
		}
		br.execTask()
	}
	return
}

func (br *BenchClient) execTask() {

	task, ok, err := br.getNextTask()

	if err != nil {
		br.log.Println("Task request failed. Details: ", err, ". Sleep 5s")
		time.Sleep(5 * time.Second)
		return
	}

	if !ok {
		br.log.Println("No task assigned. Sleep 2s")
		time.Sleep(2 * time.Second)
		return
	}

	br.log.Println("Next task to fullfil: Benchmark ", task.PackageUrl)
	result := common.TaskResult{Id: task.Id, Round: make(map[string]parse.Set)}

	// download target package
	fPath, err := downloadPackage(task.PackageUrl)
	if err != nil {
		br.log.Printf("Package download failed. Details: %s", err)
		return
	}
	br.log.Println("Package downloaded")

	// download target package dependencies
	err = downloadPackageDependencies(fPath)
	if err != nil {
		br.log.Printf("Package dependencies download failed. Details: %s", err)
		return
	}
	br.log.Println("Package dependecies downloaded")

	os.Exit(0)

	// 2. прогоняем go test bench для разного количества GOMAXPROCSs
	for i := 0; i < runtime.NumCPU(); i++ {
	    idx := "cpu" + strconv.Itoa(i)

	    // 3. Вызываем тест и парсим ответ
	    result.Round[idx], err = runTest( task.PackageUrl, i )
	    
	    if err != nil {
		log.Printf ("Failed to run test for ", task.PackageUrl, " on " , i , " CPU(s): ", err )
		continue
	    }
	    
	}

	// several attepmts to submit task execution results
	for i := 0; i < maxSubmitAttepmts; i++ {
		br.log.Printf("Result submit attempt: %d", i+1)
		clean, err := br.submitResult(&result)
		if err != nil {
			br.log.Printf("Result submit failed. Details: %s", err)
			if clean {
				break
			}
			time.Sleep(2 * time.Second)
			continue
		}
		br.log.Println("Result submited sucessfully")
		break
	}

	return
}

// getNextTask retrives next benchmarking task from goben.ch server
func (br *BenchClient) getNextTask() (*common.TaskResponse, bool, error) {

	log.Println("getNextTask started")

	buf, err := json.Marshal(common.TaskRequest{AuthKey: br.authKey, Email: br.email})
	if err != nil {
		return nil, false, err
	}

	resp, err := br.client.Post(br.nextTaskUrl, "application/json; charset=UTF-8", bytes.NewReader(buf))
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	// process sucessful response
	if resp.StatusCode == http.StatusOK {
		dec := json.NewDecoder(resp.Body)

		var task common.TaskResponse

		if err = dec.Decode(&task); err != nil {
			return nil, false, err
		}

    		log.Println( "getNextTask done. Task:", task.PackageUrl )

		return &task, len(task.PackageUrl) > 0, nil
	}

	// there is no package to process
	if resp.StatusCode == http.StatusNoContent {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, false, nil
	}

	// other server responses
	buf, _ = ioutil.ReadAll(resp.Body)
	log.Println( "getNextTask done. No task received" )
	return nil, false, errors.New(string(buf))
}

// submitResult sends task result to goben.ch server
func (br *BenchClient) submitResult(result *common.TaskResult) (bool, error) {

	buf, err := json.Marshal(result)
	if err != nil {
		return true, err
	}

	resp, err := br.client.Post(br.submitResultUrl, "application/json; charset=UTF-8", bytes.NewReader(buf))
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	// process sucessful response
	if resp.StatusCode == http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		return true, nil
	}

	// unknown tasks id or authKey, task is already done
	clean := resp.StatusCode == http.StatusBadRequest

	// process unexpeted error
	buf, _ = ioutil.ReadAll(resp.Body)
	return clean, errors.New(string(buf))
}

func (br *BenchClient) stop() {
	close(br.stopCh)
}
