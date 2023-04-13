package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	c := loadConfig(os.Args)
	log.Printf("starting mkwiki on %s", c.getServerUri())

	// TODO: add static and image handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := newPage(r.URL.Path, c)
		w.WriteHeader(p.responseCode)
		if err := c.layout.Execute(w, p); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	})
	if err := http.ListenAndServe(c.getServerUri(), nil); err != nil {
		log.Fatal(err)
	}
}
