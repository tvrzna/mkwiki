package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tvrzna/go-utils/args"
)

var buildVersion string

type theme byte

const (
	themeLight theme = iota
	themeDark
)

func (t *theme) getStyle() string {
	switch *t {
	case themeDark:
		return "theme-dark.css"
	case themeLight:
	default:
		return "theme-light.css"
	}
	return "theme-light.css"
}

type config struct {
	path   string
	appUrl string
	port   int
	layout *template.Template
	theme
}

func loadConfig(arg []string) *config {
	cwd, _ := os.Getwd()

	c := &config{cwd, "", 1500, layout(), themeLight}

	args.ParseArgs(arg, func(arg, value string) {
		switch arg {
		case "-h", "--help":
			fmt.Printf(`Usage: mkwiki [options]
Options:
	-h, --help			print this help
	-v, --version			print version
	-t, --path [PATH]		absolute path to markdown storage
	-p, --port [PORT]		sets port for listening
	-a, --app-url [APP_URL]		application url (if behind proxy)
	-e, --theme [light|dark]	sets color theme
`)
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("mkwiki %s\nhttps://github.com/tvrzna/mkwiki\n\nReleased under the MIT License.\n", c.getVersion())
			os.Exit(0)
		case "-t", "--path":
			if path, err := filepath.Abs(value); err != nil {
				log.Fatal("wrong path", err)
			} else {
				c.path = path
			}
		case "-p", "--port":
			c.port, _ = strconv.Atoi(value)
		case "-a", "--app-url":
			c.appUrl = value
		case "-e", "--theme":
			c.theme = parseTheme(value)
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

func (c *config) getVersion() string {
	if buildVersion == "" {
		return "develop"
	}
	return buildVersion
}

func parseTheme(val string) theme {
	switch strings.ToLower(val) {
	case "dark":
		return themeDark
	case "light":
	default:
		return themeLight
	}
	return themeLight
}
