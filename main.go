package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	server := http.NewServeMux()
	c := loadConfig(os.Args)
	log.Printf("starting mkwiki on %s", c.getServerUri())

	server.HandleFunc("/", requestHandler(c))
	if err := http.ListenAndServe(c.getServerUri(), server); err != nil {
		log.Fatal(err)
	}
}
