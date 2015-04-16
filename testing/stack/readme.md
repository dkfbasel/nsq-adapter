# STACK TEST
This test will provide you with a dummy stack to test the nsq throughput. You need to have a running instance of a nsqd- and a nsqlookupd-service to run the test.

## USAGE
- Start a nsqlookupd and a connected nsqd-service (I recommend respective docker images)
- Set the nsqlookupd-address constant in the repository-service and web-service
- Set the address of the http-server to be used for the web-service in the web-service and the client-service
- Set the testing parameters in the client-service
- Compile the code
- Start the repository-service
- Start the web-service
- Run the client-service
