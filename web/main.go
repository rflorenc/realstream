package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":8081", "website address")
	flag.Parse()
	mux := http.NewServeMux()

	// serve static files from a folder called "public"
	mux.Handle("/", http.StripPrefix("/",
		http.FileServer(http.Dir("public"))))

	log.Println("Serving website at:", *addr)
	http.ListenAndServe(*addr, mux)
}
