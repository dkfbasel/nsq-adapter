package nsqAdapter

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bitly/go-nsq"
)

// NewNsqAdapter will create a new nsq-adapter using the given address to connect
// to a nsqlookupd-service
func New(serviceName string, nsqlookupHttpAddress string) *NsqAdapter {

	// initialize a new adapter
	queue := NsqAdapter{
		Name:             serviceName,
		nsqlookupAddress: nsqlookupHttpAddress,
		consumers:        make(map[string]*nsq.Consumer),
		requests:         make(map[string]chan Message),
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
	message.Payload = payload

	// set the type
	message.MessageType = messageType

	return &message
}
