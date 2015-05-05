package nsqAdapter

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	responseSuffix      string = "-responses"
	responseChannelName string = "responses"
)

// InitResponseHandling will register a consumer on a topic called [NAME]-responses
func (queue *NsqAdapter) InitializeResponseHandling() {

	// create a new channel that receives all responses
	responses := make(chan Message)

	// publish an heartbeat message to the queue to handle request promptly from the start
	queue.Publish(queue.Name+responseSuffix, "init")

	// create a new consumer to handle responses for requests from this service
	queue.Subscribe(queue.Name+responseSuffix, responseChannelName, responses)

	// handle the responses in a separate go-routine
	go func() {

		// we are going to need a mutex to look access to our requests queue
		var mutex = &sync.RWMutex{}

		for {
			// wait for messages to arrive on the response-topic
			message := <-responses

			// only responses should be handled
			// note: everything else will just be discarded (i.e. init messages)
			if message.MessageType == MessageTypeResponse {
				// prevent synchronous access to the queue
				mutex.RLock()
				forwardChannel, exists := queue.requests[message.Request.OriginId]
				mutex.RUnlock()

				if exists == false {
					// the request was not found in the in-memory request queue
					continue
				}

				// send the message received to the channel specified
				// in the request-queue map. note: the SendRequest-method
				// will take care of removing the request from the queue
				// and closing the channel
				forwardChannel <- message
			}
		}
	}()
}

// SendRequest will send the given message to the given topic and wait for a response
func (queue *NsqAdapter) SendRequest(topic string, payload interface{}, timeout time.Duration) (*Message, error) {

	// create a mutex to access our map of requests
	var mutex = &sync.RWMutex{}

	// create a new Message
	message := queue.NewMessage(topic, MessageTypeRequest, payload)

	// set a timeout on the message
	message.Request.Until = time.Now().Add(timeout).String()

	// define the name of the response topic
	message.Request.RespondTo = queue.Name + responseSuffix

	// create a new response channel
	responseChannel := make(chan Message)

	// add the request to our request-queue
	// (add mutex locks to protect access from multiple go routines)
	mutex.Lock()
	queue.requests[message.Id] = responseChannel
	mutex.Unlock()

	defer func() {
		// remove the request from our request-queue
		// (add mutex lock to protect access from multiple go routines)
		mutex.Lock()
		delete(queue.requests, message.Id)
		mutex.Unlock()

		// close the request channel
		close(responseChannel)
	}()

	// publish our message to the specified topic
	err := queue.PublishMessage(topic, message)
	if err != nil {
		return nil, err
	}

	// wait for response our timeout
	select {
	case response := <-responseChannel:
		// wait for a response on the response channel
		fmt.Println("received response")
		return &response, nil

	case <-time.After(timeout):
		// wait for timeout
		err = errors.New("Request timed out:" + timeout.String())
		return nil, err
	}

}

// RespondTo will send a response to a message sent by a client via SendRequest
func (queue *NsqAdapter) RespondTo(message Message, responsePayload interface{}) error {

	// create a new response message
	responseMessage := queue.NewMessage(message.Request.RespondTo, MessageTypeResponse, responsePayload)

	// set the end time of the response
	responseMessage.EndTime = time.Now().String()

	// set the originating request id
	responseMessage.Request.OriginId = message.Id

	return queue.PublishMessage(message.Request.RespondTo, responseMessage)
}

// ForwardRequest will forward the given message to a different topic, the message payload can be appended or overwritten
func (queue *NsqAdapter) ForwardRequest(topic string, message *Message) error {
	// publish our message to the specified topic
	return queue.PublishMessage(topic, message)
}
