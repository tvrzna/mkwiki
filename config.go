package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/tvrzna/go-utils/args"
)

type config struct {
	path   string
	appUrl string
	port   int
	layout *template.Template
}

func loadConfig(arg []string) *config {
	cwd, _ := os.Getwd()

	c := &config{cwd, "", 1500, layout()}

	args.ParseArgs(arg, func(arg, value string) {
	args.ParseArgs(os.Args, func(arg, value string) {
		switch arg {
		case "-h", "--help":
			fmt.Printf("Usage: mkwiki [options]\nOptions:\n\t-h, --help\t\t\tprint this help\n\t-v, --version\t\t\tprint version\n\t-t, --path [PATH]\t\tabsolute path to markdown storage\n\t-p, --port [PORT]\t\tsets port for listening\n\t-a, --app-url [APP_URL]\t\tapplication url (if behind proxy)\n")
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("mkwiki 0.1.0\nhttps://github.com/tvrzna/mkwiki\n\nReleased under the MIT License.\n")
			os.Exit(0)
		case "-t", "--path":
			c.path = value
		case "-p", "--port":
			c.port, _ = strconv.Atoi(value)
		case "-a", "--app-url":
			c.appUrl = value
		}
	})

	return c
}

func (c *config) getServerUri() string {
	return "127.0.0.1:" + strconv.Itoa(c.port)
}

func (c *config) getAppUrl() string {
	if c.appUrl == "" {
		return "http://" + c.getServerUri()
	}
	return c.appUrl
}
