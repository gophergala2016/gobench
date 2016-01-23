package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/tools/benchmark/parse"
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

const baseUrl = "https://www.magnova.ru"

type TaskRequest struct {
	AuthKey string `json: "authKey"`
	Email   string `json: "email"`
}

type TaskResponse struct {
	PackageUrl string `json:"packageUrl"`
}

type TaskResult struct {
	Config
	AuthKey       string `json: "authKey"`
	Email         string `json: "email"`
	Specification string `json: "specification"`

	// Result holds parsed bencmark results per GoMaxProc 1-8
	Result  map[string]parse.Benchmark
	BuildOk bool `json:"buildOk"`
}

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

	packageName, ok, err := br.getNextTask()
	if err != nil {
		br.log.Println("Task request failed. Details: ", err)
		time.Sleep(20 * time.Second)
		return
	}
	if !ok {
		br.log.Println("No tasks to do. Sleep 10s")
		time.Sleep(10 * time.Second)
		return
	}

	br.log.Println("Next task to do: ", packageName)

	// TODO:
	// 1. выкачиваем пакеты и зависимости
	// 2. прогоняем go test bench и т.д
	// 3. парсим ответ
	// 4. отправляем на сервер, вместе с параметрами тестового окружения

	return
}

// nextTask retrives next benchmarking task from goben.ch server
func (br *BenchRunner) getNextTask() (string, bool, error) {

	taskReq := TaskRequest{AuthKey: br.authKey, Email: br.email}
	buf, err := json.Marshal(taskReq)
	if err != nil {
		return "", true, err
	}

	resp, err := br.client.Post(baseUrl+"/shop/api/signIn", "application/json; charset=UTF-8", bytes.NewReader(buf))
	if err != nil {
		return "", true, err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO: handle errors
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		return "", true, errors.New(resp.Status)
	}

	dec := json.NewDecoder(resp.Body)

	var taskResp TaskResponse

	if err = dec.Decode(&taskResp); err != nil {
		resp.Body.Close()
		return "", true, errors.New(resp.Status)
	}

	resp.Body.Close()
	return taskResp.PackageUrl, true, nil
}

func (br *BenchRunner) stop() {
	close(br.stopCh)
}
