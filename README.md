# nsq-adapter
A small adapter library that makes using the awesome nsq-service a snap

# usage
```go
package main

import (
	"github.com/tikiatua/nsq-adapter"
)

const (
	nsqlookupdAddress string = "192.168.99.100:4161"
)

func main() {

	// create a new nsqadapter with the name of our service
	queue := nsqAdapter.New("name-of-my-service", nsqlookupdAddress)
	
	// initialize the ability to handle responses.
	// this will create a topic called [name-of-service]-responses
	// that will be used to receive responses to requests issued
	// from this service
	queue.InitializeResponseHandling()

	// create a channel that will receive message from
	// a topic we would like to subscribe to
	messageChan := make(chan nsqAdapter.Message)
	
	// subscribe to any topic and channel
	// (note: multiple subscribers to the same channel can be used for load balancing)
	queue.Subscribe("please-do-something-with-this", "requests", messageChan)

	// handle all incoming messages of our subscription
	for {
  		// wait for incoming messages
		message := <-messageChan
	
		// handle the messages
		fmt.Println("RECEIVED:", message.Payload)

    	// send a response if the message is a request
		if message.MessageType == nsqAdapter.MessageTypeRequest {
			queue.RespondTo(message, "this is an answer from name-of-my-service")
		}
	}
}
```
You can find an example usage with a testclient, a webserver and a simple repository service in the testing directory
