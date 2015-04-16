package nsqAdapter

import (
	"encoding/json"

	"github.com/bitly/go-nsq"
)

// Subscribe will start listening on the given topic for messages
func (queue *NsqAdapter) Subscribe(topic string, channel string, messageChannel chan<- Message) error {

	var err error

	// combine topic and channel as address
	domain := topic + "." + channel

	// check if we do already have a consumer for this channel
	consumer, ok := queue.consumers[domain]

	// create a new consumer if necessary
	if ok == false {
		// create a new non concurrent consumer
		consumer, err = nsqConsumer(topic, channel, queue.config)
		if err != nil {
			return err
		}

		// remember the consumer for later
		queue.consumers[domain] = consumer
	}

	// create a new handler function that will parse all nsq-messages
	// and submit it to the supplied message channel
	nsqHandler := func(nsqMessage *nsq.Message) error {
		// parse the nsq-message body into our custom message struct
		message := Message{}
		err := json.Unmarshal(nsqMessage.Body, &message)
		if err != nil {
			return err
		}

		// return the data to our messageChannel
		messageChannel <- message
		return nil
	}

	// attach our handler function to the consumer
	attachHandler(consumer, nsqHandler, queue.nsqlookupAddress)

	return nil
}

// Publish will wrap the given payload into a custom message and send
// it off to nsq (without waiting for a response)
func (queue *NsqAdapter) Publish(topic string, payload interface{}) error {

	// create a new encapsulated message
	message := queue.NewMessage(topic, MessageTypeBroadcast, payload)

	// publish our predefined message
	return queue.PublishMessage(topic, message)
}

// PublishMessage will publish a message struct to nsq and forget about it
func (queue *NsqAdapter) PublishMessage(topic string, message *Message) error {

	// check if we do already have a producer for this topic
	if queue.producer == nil {
		// create a new producer
		producer, err := nsqProducer(queue.nsqlookupAddress, queue.config)
		if err != nil {
			return err
		}

		// remember the producer for later
		queue.producer = producer
	}

	// marshal the message for transport to json
	messageString, _ := json.Marshal(&message)

	// publish the message
	return queue.producer.Publish(topic, messageString)
}
