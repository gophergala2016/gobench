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
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	baseUrl = "http://127.0.0.1:8080"
	debug   = true
)

// BenchClient implements client to goben.ch server
type BenchClient struct {
	client        *http.Client
	log           *log.Logger
	authKey       string
	email         string
	stopCh        chan os.Signal
	specification string
}

// NewBenchClient creates BenchClient instance
func NewBenchClient(authKey, email string, l *log.Logger) (*BenchClient, error) {

	br := &BenchClient{
		authKey: authKey,
		email:   email,
		client:  &http.Client{Timeout: 2 * time.Second},
		log:     l,
	}

	// TODO: identify current machine specification: RAM, CPU, OS, etc.
	// save to br.specification

	err := br.Ping()
	if err != nil && debug == false {
		// if debug mode is on, exit with error
		return nil, err
	}
	return br, nil
}

// Ping checks server availability
func (br *BenchClient) Ping() error {

	resp, err := br.client.Head(baseUrl)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
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

	// Download target package
	cmd := exec.Command("go", "get", task.PackageUrl)
	err = cmd.Start()
	if err != nil {
		br.log.Printf("Package download failed. Details: %s", err)
		// TODO: inform server abour problem
		return
	}

	log.Printf("Waiting for command to finish...")
	if err = cmd.Wait(); err != nil {
		br.log.Printf("Command finished with error: %s", err)
		// TODO: проверить как это работает на самом деле
	}

	// Donwload dependencies
	// TODO

	//
	//br.log.Println(err, string(buf))
	os.Exit(0)

	// Отсюда и ниже уже

	// 2. прогоняем go test bench для разного количества GOMAXPROCSs
	for i := 0; i < runtime.NumCPU(); i++ {
		// 3. парсим ответ

	}

	// 4. отправляем на сервер, вместе с параметрами тестового окружения

	err = br.submitResult(&result)
	if err != nil {
		br.log.Println("Result submit failed")
		return
	}

	br.log.Println("Result submited sucessfully")

	return
}

// getNextTask retrives next benchmarking task from goben.ch server
func (br *BenchClient) getNextTask() (*common.TaskResponse, bool, error) {

	buf, err := json.Marshal(common.TaskRequest{AuthKey: br.authKey, Email: br.email})
	if err != nil {
		return nil, false, err
	}

	resp, err := br.client.Post(baseUrl+"/api/task/next", "application/json; charset=UTF-8", bytes.NewReader(buf))
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

		return &task, len(task.PackageUrl) > 0, nil
	}

	// there is no package to process
	if resp.StatusCode == http.StatusNoContent {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, false, nil
	}

	// other server responses
	buf, _ = ioutil.ReadAll(resp.Body)
	return nil, false, errors.New(string(buf))
}

// submitResult sends benchmark result to goben.ch server
func (br *BenchClient) submitResult(result *common.TaskResult) error {

	buf, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := br.client.Post(baseUrl+"/api/task/submit", "application/json; charset=UTF-8", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// process sucessful response
	if resp.StatusCode == http.StatusOK {
		dec := json.NewDecoder(resp.Body)

		var taskResp common.TaskResponse

		if err = dec.Decode(&taskResp); err != nil {
			return err
		}

		return nil
	}

	// there is no package to process
	if resp.StatusCode == http.StatusNoContent {
		io.Copy(ioutil.Discard, resp.Body)
		return nil
	}

	// unexpected error
	buf, _ = ioutil.ReadAll(resp.Body)
	return errors.New(string(buf))
}

func (br *BenchClient) stop() {
	close(br.stopCh)
}
