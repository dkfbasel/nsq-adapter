package nsqAdapter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// define responce from node
type nodesResponse struct {
	Status_code int
	Status_txt  string
	Data        producers
}

// multiple-nsq-producers
type producers struct {
	Producers []producer
}

// one nsq-producer
type producer struct {
	Remote_address    string
	Hostname          string
	Broadcast_address string
	Version           string
	Http_port         int
	Tcp_port          int
	Tombstones        []string
	Topics            []string
}

// discoverNsqd will find all producing nsqd-nodes from the specified nsqlookupd-address.
// It will return a list of nsqd addresses.
func discoverNsqd(nsqlookupdHttp string) ([]string, error) {

	// find all producing nsqd-nodes from the specified nsqlookupd-address
	resp, err := http.Get("http://" + nsqlookupdHttp + "/nodes")
	if err != nil {
		log.Fatal(err)
	}

	// read the responnse content
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// parse the response into a struct
	var n nodesResponse
	json.Unmarshal(body, &n)

	// check if status code is ok
	if n.Status_code != 200 {
		return nil, errors.New("could not fetch nsqd nodes")
	}

	// check if there are any registered producers
	numberOfProducers := len(n.Data.Producers)
	if numberOfProducers <= 0 {
		return nil, errors.New("could not find any nsqd producers")
	}

	// create a new slice for our addresses
	producerAddresses := make([]string, numberOfProducers)

	for index, producer := range n.Data.Producers {
		producerAddresses[index] = producer.Broadcast_address + ":" + strconv.Itoa(producer.Tcp_port)
	}

	return producerAddresses, nil
}
