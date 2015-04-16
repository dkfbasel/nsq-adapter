package nsqAdapter

import "github.com/bitly/go-nsq"

// handlerFunction is used to pass functions to the nsq-consumer instead
// of struct that implements a HandleMessage method
type nsqHandlerFunction func(nsqMessage *nsq.Message) error

func nsqConsumer(topic string, channel string, nsqConfig *nsq.Config) (*nsq.Consumer, error) {

	// create a new nsq-consumer
	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func attachHandler(consumer *nsq.Consumer, handler nsqHandlerFunction, nsqlookupdHttpAddress string) {
	// add a new non concurrent handler
	consumer.AddHandler(nsq.HandlerFunc(handler))

	// connect our consumer to the nsq-lookup-service
	consumer.ConnectToNSQLookupd(nsqlookupdHttpAddress)
}
