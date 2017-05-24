package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/bitly/go-nsq"

	"sync"

	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	const updateDuration = 1 * time.Second
	/*
		deferred statements are run in LIFO	order,
		the first function we defer will be the last function to be executed, which
		is why the first thing we do  here is to defer the exiting code.
	*/
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connecting to database...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	/*
		this code will run before our previously deferred statement containing
		the exit code, because deferred functions are run in the reverse order in which
		they were called. Therefore, whatever happens, we know that the
		database session will properly close.
	*/
	defer func() {
		log.Println("Closing database connection...")
		db.Close()
	}()

	pollData := db.DB("ballots").C("polls")

	// Consuming NSQ messages
	var counts map[string]int
	var countsLock sync.Mutex

	log.Println("Connecting to NSQ...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}
	log.Println("Waiting for votes on NSQ...")
	var updater *time.Timer

	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("Updating database...")
			log.Println(counts)
			ok := true
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if _, err := pollData.UpdateAll(sel, up); err != nil {
					log.Println("failed to update:", err)
					ok = false
				}
			}
			if ok {
				log.Println("Finished updating database...")
				counts = nil // reset counts
			}
		}
		updater.Reset(updateDuration)

	})
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-termChan:
			updater.Stop()
			q.Stop()
		case <-q.StopChan:
			return
		}
	}
}
