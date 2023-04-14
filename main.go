package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
)

func main() {
	server := http.NewServeMux()
	c := loadConfig(os.Args)
	log.Printf("starting mkwiki on %s", c.getServerUri())

	// TODO: add image handler
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/style.css", "/favicon.ico":
			f, _ := www.ReadFile("www" + r.URL.Path)
			w.Header().Add("content-type", mime.TypeByExtension(r.URL.Path[strings.LastIndex(r.URL.Path, "."):]))
			w.Write(f)
		default:
			p := newPage(r.URL.Path, c)
			w.WriteHeader(p.responseCode)
			if err := c.layout.Execute(w, p); err != nil {
				w.WriteHeader(500)
				log.Println(err)
			}
		}
	})
	if err := http.ListenAndServe(c.getServerUri(), server); err != nil {
		log.Fatal(err)
	}
}
