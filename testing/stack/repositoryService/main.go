package main

import (
	"github.com/kr/pretty"
	"github.com/tikiatua/nsq-adapter"
)

const (
	nsqlookupd string = "192.168.99.100:4161"
)

func main() {

	// create a new nsqadapter
	queue := nsqAdapter.New("test-arena-repository", nsqlookupd)

	// subscribe to our repository channel
	repositoryChan := make(chan nsqAdapter.Message)
	queue.Subscribe("test-arena-repository", "requests", repositoryChan)

	// handle all messages
	for {
		message := <-repositoryChan
		pretty.Println("RECEIVED:", message.Payload)

		if message.MessageType == nsqAdapter.MessageTypeRequest {
			queue.RespondTo(message, message.Payload.(string)+"-from-repository")
		}
	}
}
