# SIMPLE TEST
This test will initialize a service that will send and receive messages. You need to have a running instance of a nsqd- and a nsqlookupd-service to run the test.

## USAGE
- Start a nsqlookupd and a connected nsqd-service (I recommend respective docker images)
- Set the nsqlookupd-address constant in the source code
- Compile and run the code
- Enter your messages in the terminal

Use a point to separate the topic to post to from the message body.
```
repository.my custom message
```

Use the following syntax to issue a request
```
repository.request.my custom request
```

The service will subscribe to the topic "webserver" and respond to all
requests sent to this topic
