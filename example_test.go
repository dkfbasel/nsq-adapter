package nsqAdapter

import "fmt"

const (
	nsqlookupdAddress string = "192.168.99.100:4161"
)

func ExampleMain() {

	// create a new nsqadapter with the name of our service
	queue := New("name-of-my-service", nsqlookupdAddress)

	// initialize the ability to handle responses.
	// this will create a topic called [name-of-service]-responses
	// that will be used to receive responses to requests issued
	// from this service
	queue.InitializeResponseHandling()

	// create a channel that will receive message from
	// a topic we would like to subscribe to
	messageChan := make(chan Message)

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
		if message.MessageType == MessageTypeRequest {
			queue.RespondTo(message, "this is an answer from name-of-my-service")
		}
	}
}
