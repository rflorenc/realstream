package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	flag.Parse()
	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	// register a single handler for all requests that begin with
	// the path /polls/
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db,
		withAPIKey(handlePolls)))))

	log.Println("Starting web server on", *addr)

	// time.Duration allows any in-flight requests some time to complete
	// before the function exits.
	// The Run function will block until the program is terminated.
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("Stopping...")
}

// Cross-browser resource sharing
// consider using https://github.com/fasterness/cors
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)

	}
}

/*
withData is a HandlerFunc wrapper that manages the database session for us

The returned http.HandlerFunc type will copy the database session,
defer the closing of that copy, and set a reference to the ballots database
as the db variable using our SetVar helper, before finally calling the next HandlerFunc .

*/
func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

func isValidAPIKey(key string) bool {
	// TODO: Extend to read API key from s3, gce,
	// or from an in memory source.

	// Determine using plain config file ?
	return key == "abc123"
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "Invalid API key")
			return
		}
		fn(w, r)
	}
}
