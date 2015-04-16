package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tikiatua/nsq-adapter"
)

const (
	webserviceAddress string = "localhost:5353"
	nsqlookupd        string = "192.168.99.100:4161"
)

var queue *nsqAdapter.NsqAdapter

func main() {

	// create a new nsqadapter
	queue = nsqAdapter.New("test-arena-webservice", nsqlookupd)

	// initialize the ability to handle responses
	queue.InitializeResponseHandling()

	// init a webserver
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// register a handler for our one and only route
	router.GET("/:id", requestHandler)

	// run the webserver
	router.Run(webserviceAddress)
}

// handle client requests (each in its own goroutine)
func requestHandler(c *gin.Context) {
	id := c.Params.ByName("id")

	// request the data from the repository service (allow maximum timeout of 5 seconds)
	message, err := queue.SendRequest("test-arena-repository", id, time.Second*5)

	if err != nil {
		fmt.Println("ERROR:", err.Error())
		c.String(500, "could not fetch data from the repository service:\n\n"+err.Error())
		return
	}

	result := message.Payload.(string)
	c.String(200, result)
}
