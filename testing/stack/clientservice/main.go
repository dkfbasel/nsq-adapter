package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

const (
	webserviceAddress                string = "http://localhost:5353"
	numberOfRequests                 int    = 500
	millisecondsDelayBetweenRequests        = 2
	verbose                          bool   = true
)

var requests []*Request

type Request struct {
	Id       string
	Url      string
	Status   string
	Started  time.Time
	Ended    time.Time
	Duration time.Duration
	Info     string
}

func (request *Request) String() string {
	var stringRep string = ""

	//stringRep = "Id:\t\t" + request.Id
	stringRep += "\nStatus:\t\t" + request.Status
	stringRep += "\nDuration:\t" + request.Duration.String()
	stringRep += "\nInfo:\t\t" + request.Info
	return stringRep
}

const (
	StatusStarted  = "started"
	StatusError    = "error"
	StatusFinished = "finished"
	StatusTimeout  = "timeout"
)

func main() {

	fmt.Println("Starting Requests ..\n")

	// initialize our request map
	requests = []*Request{}

	var waitGroup sync.WaitGroup
	// done := make(chan struct{})

	// issue a http request
	for i := 0; i < numberOfRequests; i++ {
		waitGroup.Add(1)
		<-time.After(millisecondsDelayBetweenRequests * time.Millisecond)
		go issueHttpRequest(webserviceAddress, 2000, &waitGroup)
	}

	// block execution
	waitGroup.Wait()

	// calculate mean duration of request
	sum := 0 * time.Millisecond
	minimum := 20000 * time.Millisecond
	maximum := 0 * time.Millisecond
	finished := 0
	timeouts := 0
	errors := 0
	unknown := 0

	for _, request := range requests {
		sum += request.Duration

		if minimum > request.Duration {
			minimum = request.Duration
		}

		if maximum < request.Duration {
			maximum = request.Duration
		}

		switch request.Status {
		case "finished":
			finished += 1
		case "error":
			errors += 1
		case "timeout":
			timeouts += 1
		default:
			unknown += 1
		}
	}

	// calculate mean duration
	meanNanoSeconds := sum.Nanoseconds() / int64(len(requests))
	meanDuration := time.Duration(meanNanoSeconds)

	if verbose {
		fmt.Println("\n- - - - - - - - - - - - - - - - - - - - - - - - - -")
	}

	fmt.Println("No. Requests:\t", len(requests))
	fmt.Println("No. Finished:\t", finished)
	fmt.Println("No. Timeout:\t", timeouts)
	fmt.Println("No. Errors:\t", errors)
	fmt.Println("No. Unknown:\t", unknown)

	fmt.Println("\nMinimum:\t", minimum.String())
	fmt.Println("Maximum:\t", maximum.String())
	fmt.Println("Mean:\t\t", meanDuration.String())

	fmt.Println("\nFinished all requests ..\n")
}

func issueHttpRequest(address string, timeoutInMilliseconds time.Duration, waitGroup *sync.WaitGroup) {

	// create a new request
	request := &Request{}
	request.Id = uuid.NewUUID().String()
	request.Url = address + "/" + request.Id
	request.Started = time.Now()

	responseChan := make(chan *Request)

	// run our request in a new go routine
	go func(request Request) {
		response, err := http.Get(request.Url)
		if err != nil {
			request.Status = StatusError
			request.Info = "Response error: " + err.Error()
			responseChan <- &request
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			request.Status = StatusError
			request.Info = "Read response error: " + err.Error()
			responseChan <- &request
			return
		}
		request.Info = string(body)
		request.Status = StatusFinished
		responseChan <- &request
		return
	}(*request)

	select {
	case request = <-responseChan:
		request.Ended = time.Now()
		request.Duration = request.Ended.Sub(request.Started)

	case <-time.After(timeoutInMilliseconds * time.Millisecond):
		request.Status = StatusTimeout
		request.Info = "--"
		request.Ended = time.Now()
		request.Duration = request.Ended.Sub(request.Started)
	}

	var mutex = &sync.Mutex{}

	mutex.Lock()
	if verbose {
		fmt.Println(request.String())
	}
	requests = append(requests, request)
	mutex.Unlock()

	waitGroup.Done()

}
