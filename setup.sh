#!/bin/bash
source twitter.env

nsqlookupd &
nsqd --lookupd-tcp-address=localhost:4160 &

if [[ ! -d db ]]; then
    mkdir -vp db
else
    mongod --dbpath ./db &
fi

for app in counter twittervotes api web
do
    cd $app
    go get -v
    go build -v -o $app
    ./$app &
    cd -
done

echo "Open a browser and head to http://localhost:8081/"

echo "Done"
exit 0
