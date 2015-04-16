package nsqAdapter

import (
	"errors"

	"github.com/bitly/go-nsq"
)

// nsqProducer will create a new producer connected to an nsqd service
func nsqProducer(nsqlookupdHttpAddress string) (*nsq.Producer, error) {

	// check with nsqlookupd if there are any nsqd-hosts
	producers, errDiscover := discoverNsqd(nsqlookupdHttpAddress)
	if errDiscover != nil {
		return nil, errDiscover
	}

	// we are going to need an nsqd host first
	if len(producers) == 0 {
		return nil, errors.New("No active nsqd hosts found")
	}

	// create a new producer to send messages
	// note: this will publish to the first producers in the given list
	producer, err := nsq.NewProducer(producers[0], nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	return producer, nil
}
