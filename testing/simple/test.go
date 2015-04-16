package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/tikiatua/nsq-adapter"
)

const (
	nsqlookupd = "192.168.99.100:4161"
)

func main() {

	// create a channel for message communication
	messages := make(chan string)

	// wait channel
	done := make(chan bool)

	// parse console input
	go parseCommandline(messages, done)

	// simulate web-server
	go testAdapter(messages, done)

	<-done

	fmt.Println("System shut down")
}

// parseMessages will parse the command line and send
// the input to a channel
func parseCommandline(messages chan<- string, done chan<- bool) {

	// use the bufio-scanner to read command line with spaces
	scanner := bufio.NewScanner(os.Stdin)
	var message string

	// keep on scanning
	for scanner.Scan() {
		message = scanner.Text()
		if message == "quit" {
			// explore option to close channel with close(done)
			done <- true
		}
		messages <- message
	}
}

func testAdapter(messages <-chan string, done <-chan bool) {

	// create a new nsqadapter
	queue := nsqAdapter.New("test", nsqlookupd)

	// initialize the ability to handle responses
	queue.InitializeResponseHandling()

	// subscribe to a certain topic
	webserverChan := make(chan nsqAdapter.Message)
	queue.Subscribe("webserver", "requests", webserverChan)

	for {
		select {
		case info := <-webserverChan:
			pretty.Println("WEBSERVER:", info.Payload)

			if info.MessageType == nsqAdapter.MessageTypeRequest {
				queue.RespondTo(info, "this is a response")
			}

		case message := <-messages:
			data := strings.Split(message, ".")

			if data[1] == "request" {

				go func() {
					pretty.Println("REQUEST:", data[0], data[2])
					// create a request  a request
					result, err := queue.SendRequest(data[0], data[2], time.Second*10)

					if err != nil {
						pretty.Println("RESPONSE:", err.Error())
					} else {
						pretty.Println("RESPONSE:", result)
					}
				}()

			} else {
				pretty.Println("PUBLISH:", data[0], data[1])
				queue.Publish(data[0], data[1])
			}

		case <-done:
			fmt.Println("STOPPING QUEUE")
		}
	}
}
