# TO RUN THE APPLICATION, YOU ARE GONNA NEED THE FOLLOWING IMAGES
nsqio/nsq:latest

# START DOCKER CONTAINERS IN THE SPECIFIED ORDER WITH THE FOLLOWING COMMANDS
# (THIS ASSUMES YOUR DOCKER HOST HAS IP 192.168.99.100)

# run nsqlookupd service to find nsqd instances
docker run --name nsqlookupd -d --publish 4160:4160 --publish 4161:4161 nsqio/nsq /nsqlookupd

# run nsqd service for message queues (map local directory into docker container for message storage)
# docker run --name nsqd -d --publish 4150:4150 --publish 4151:4151 --volume /Users/tiki/gocode/src/github.com/tikiatua/wissenschaftspool/data/nsq:/data dkfbasel/nsqd --broadcast-address=192.168.59.103 --lookupd-tcp-address=192.168.59.103:4160 --data-path=/data

docker run --name nsqd -d --publish 4150:4150 --publish 4151:4151 nsqio/nsq /nsqd--broadcast-address=192.168.99.100 --lookupd-tcp-address=192.168.99.100:4160

# run nsqadmin service for nsqd administration through web-browser (on port 4171)
docker run --name nsqadmin -d --publish 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=192.168.99.100:4161
