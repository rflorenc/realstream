#!/bin/bash

source twitter.env

nsqlookupd
nsqd --lookupd-tcp-address=localhost:4160

mongod --dbpath ./db &
# > use ballots
# switched to db ballots
# > db.polls.insert({"title":"Test poll","options":["one","two","three"]})
# > db.polls.insert({"title":"Test poll two","options":["four","five","six"]})
#

./counter
./twittervotes
./api
# curl -s  http://localhost:8080/polls/?key=abc123 | python -m json.tool
./web


#nsq_tail --topic="votes" --lookupd-http-address=localhost:4161


