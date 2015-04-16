package nsqAdapter

import "github.com/bitly/go-nsq"

// NsqAdapter is our main struct for external packages to interact with
type NsqAdapter struct {
	Name             string      // then name of this service
	config           *nsq.Config // the configuration for the nsq service
	nsqlookupAddress string      // the address of the nsqlookup-service

	producer  *nsq.Producer            // a producer to send messages
	consumers map[string]*nsq.Consumer // a map of all our consumers

	requests        map[string]chan Message // a map of all pending requests
	responseHandler *nsq.Consumer           // handles all incoming responses to our service
}

// Message is a struct to hold all messages sent over nsq
type Message struct {
	Id        string      // a unique id for the message
	From      string      // the originating service
	To        string      // the topic that the message is posted to
	StartTime string      // the time the processing was started
	EndTime   string      // the time the process has ended
	Payload   interface{} // the payload of the data

	MessageType string  // the type of message (i.e. broadcast, request, response)
	Request     Request // specific information for requests
}

// Request is a struct to hold additional information for a request message
type Request struct {
	Until     string // until which moment does the request have to be fulfilled
	RespondTo string // the topic to respond to
	OriginId  string // the id of the originating response
}

const (
	MessageTypeRequest   string = "request"
	MessageTypeResponse  string = "response"
	MessageTypeBroadcast string = "broadcast"
)
