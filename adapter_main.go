// The package nsqAdapter provides a thin wrapper around bitly's nsq-services to
// simplify usage of nsq's asynchronous message queues for midsized projects.
//
// Basically, you define a new nsq-adapter struct and call the following methods:
// - Subscribe() to subscribe to a specific topic and handle incoming messages
// - Publish() to send a message to a specific topic
// - Request() to send a request to a specific topic and wait for a response
// - RespondTo() to send a response message to a request
//
// Please take a look at the example and the samples in the testing directory for
// an idea about the usage of the package.
package nsqAdapter

import (
	"encoding/json"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bitly/go-nsq"
)

// NewNsqAdapter will create a new nsq-adapter using the given address to connect
// to a nsqlookupd-service and use the default configuration for connections
func New(serviceName string, nsqlookupHttpAddress string) *NsqAdapter {
	return NewWithCustomConfig(serviceName, nsqlookupHttpAddress, nsq.NewConfig())
}

// NewWithCustomConfig will create a new nsq-adapter with a custom nsq-configuration
func NewWithCustomConfig(serviceName string, nsqlookupHttpAddress string, config *nsq.Config) *NsqAdapter {

	// initialize a new adapter
	queue := NsqAdapter{
		Name:             serviceName,
		nsqlookupAddress: nsqlookupHttpAddress,
		consumers:        make(map[string]*nsq.Consumer),
		requests:         make(map[string]chan Message),
		config:           config,
	}
	return &queue
}

// New Message will create a new message struct to send to nsq
func (queue *NsqAdapter) NewMessage(topic string, messageType string, payload interface{}) *Message {

	// create a new Message
	message := Message{}

	// set a unique id for our message
	message.Id = uuid.NewUUID().String()

	// set the originating service
	message.From = queue.Name

	// set the message to send the data to
	message.To = topic

	// define the time until we need the response
	message.StartTime = time.Now().String()

	// set the payload
	message.Payload, _ = json.Marshal(payload)

	// set the type
	message.MessageType = messageType

	return &message
}
