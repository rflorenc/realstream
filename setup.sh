#!/usr/bin/env bash
#
# TODO: 
# 1. Describe using kubernetes spec files. 
# 2. Deploy against a local minikube cluster.

# This is 80's horrible. 

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo $BASEDIR
trap sig_handler SIGINT


function sig_handler() {
    echo -e "\nReceived terminate signal. Exiting."
    cleanup
    exit 1
}

function cleanup() {
    cd $BASEDIR
    rm -rf db &> /dev/null
    rm  *.dat &> /dev/null
    exit ?
}

source twitter.env
SLEEP=5

nsqlookupd &
sleep $SLEEP

nsqd --lookupd-tcp-address=localhost:4160 &
sleep $SLEEP

mkdir -vp db &> /dev/null

mongod --dbpath ./db &
sleep $SLEEP

for app in counter twittervotes api web
do
    cd $app
    go get -v 
    go build -v -o $app
    ./$app &
    cd -
done

sleep $SLEEP

echo "Open a web browser and head to http://localhost:8081/"

echo "Done"
exit 0
