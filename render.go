package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed www
var www embed.FS

type page struct {
	responseCode int
	c            *config
	Path         string
	Content      template.HTML
	LastModify   time.Time
	ContentList  []*pageContent
}

type pageContent struct {
	Current bool
	Level   int
	Path    string
}

func newPage(path string, c *config) *page {
	p := &page{responseCode: 200, c: c}

	if path == "/" {
		path = "/index"
	}
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	filePath := c.path + "/" + path + ".md"
	filePath = strings.ReplaceAll(filePath, "//", "/")

	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	p.Path = path
	p.loadMarkdown(filePath)
	p.loadContentList(c, path)

	return p
}

func layout() *template.Template {
	tpl, err := template.ParseFS(www, "www/template.html")
	if err != nil {
		log.Fatal(err)
	}
	return tpl
}

func requestHandler(c *config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/style.css" || r.URL.Path == "/favicon.ico" {
			f, _ := www.ReadFile("www" + r.URL.Path)
			w.Header().Add("content-type", getMimeType(r.URL.Path))
			w.Write(f)
		} else if b, _ := regexp.MatchString("(?i).*\\.(png|jpg|jpeg|gif|ico)$", r.URL.Path); b {
			serveImage(c, r.URL.Path, w)
		} else {
			p := newPage(r.URL.Path, c)
			w.WriteHeader(p.responseCode)
			if err := c.layout.Execute(w, p); err != nil {
				w.WriteHeader(500)
				log.Println(err)
			}
		}
	}
}

func serveImage(c *config, path string, w http.ResponseWriter) {
	f, err := os.OpenFile(c.path+"/"+path, os.O_RDONLY, 0600)
	w.Header().Add("content-type", getMimeType(path))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(buf[:n])
	}
}

func getMimeType(path string) string {
	return mime.TypeByExtension(path[strings.LastIndex(path, "."):])
}

func (p *page) loadMarkdown(path string) {
	md, err := os.ReadFile(path)
	if err != nil {
		p.responseCode = 404
		p.Content = template.HTML("<h1>404 - not found</h1>")
		return
	}
	if fi, err := os.Stat(path); err == nil {
		p.LastModify = fi.ModTime()
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	pr := parser.NewWithExtensions(extensions)
	doc := pr.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	p.Content = template.HTML(markdown.Render(doc, renderer))
}

func (p *page) loadContentList(c *config, currentPath string) {
	filepath.WalkDir(c.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		url := strings.ReplaceAll(strings.ReplaceAll(path, c.path+"/", ""), ".md", "")

		p.ContentList = append(p.ContentList, &pageContent{url == currentPath, strings.Count(url, string("/")), url})

		return nil
	})

	sort.Slice(p.ContentList[:], func(i, j int) bool {
		if p.ContentList[i].Level == p.ContentList[j].Level {
			return p.ContentList[i].Path < p.ContentList[j].Path
		}
		return p.ContentList[i].Level < p.ContentList[j].Level
	})
}

func (p *page) UrlFor(path string) string {
	return p.c.getAppUrl() + "/" + path
}
