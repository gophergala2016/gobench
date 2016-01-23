package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gophergala2016/gobench/common"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const baseUrl = "http://127.0.0.1:8080"

// BenchRunner represents goben.ch client attributes
type BenchRunner struct {
	client  *http.Client
	log     *log.Logger
	stopCh  chan os.Signal
	authKey string `json: "authKey"`
	email   string `json: "email"`
}

// NewBenchRunner creates BenchRunner instance
func NewBenchRunner(authKey, email string, l *log.Logger) (*BenchRunner, error) {

	// TODO: определяем параметры тестового окружения
	br := &BenchRunner{authKey: authKey, email: email,
		client: &http.Client{Timeout: 2 * time.Second},
		log:    l}

	err := br.Ping()
	if err != nil {
		return nil, err
	}
	return br, nil
}

// Ping checks server availability
func (br *BenchRunner) Ping() error {

	resp, err := br.client.Head(baseUrl)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// Run starts goben.ch client
func (br *BenchRunner) Run() {

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

func (br *BenchRunner) run() {

	for {
		select {
		case <-br.stopCh:
			return
		default:
			//If the channel is still open, continue as normal
		}

		br.exec()
	}

	return
}

func (br *BenchRunner) exec() {

	task, ok, err := br.getNextTask()
	if err != nil {
		br.log.Println("Task request failed. Details: ", err, ". Sleep 5s")
		time.Sleep(5 * time.Second)
		return
	}

	if !ok {
		br.log.Println("No packages to benchmark. Sleep 2s")
		time.Sleep(2 * time.Second)
		return
	}

	br.log.Println("Next package to bench: ", task.PackageUrl)

	// TODO:
	// 1. выкачиваем пакеты и зависимости
	// 2. прогоняем go test bench и т.д
	// 3. парсим ответ
	// 4. отправляем на сервер, вместе с параметрами тестового окружения

	result := common.TaskResult{Id: task.Id}
	err = br.submitResult(&result)
	if err != nil {
		br.log.Println("Result submit failed")
		return
	}

	br.log.Println("Result submited sucessfully")

	return
}

// getNextTask retrives next benchmarking task from goben.ch server
func (br *BenchRunner) getNextTask() (*common.TaskResponse, bool, error) {

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
func (br *BenchRunner) submitResult(result *common.TaskResult) error {

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

func (br *BenchRunner) stop() {
	close(br.stopCh)
}
