package main

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type page struct {
	responseCode int
	Path         string
	Content      template.HTML
	ContentList  []*pageContent
}

type pageContent struct {
	Current bool
	Level   int
	Path    string
}

func newPage(path string, c *config) *page {
	p := &page{responseCode: 200}

	if path == "/" {
		path = "/index"
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
	layout := `<!DOCTYPE html>
	<html>
	<head>
		<title>mkwiki | {{ .Path }}</title>
	</head>
	<body>
		<div>
			<div>
				{{ .Content }}
			</div>
			<div>
				<ul>
					{{ range .ContentList }}
						<li>
							{{ if .Current }}
								{{ .Path }}
							{{ else }}
								<a href="/{{ .Path }}">{{ .Path }}</a>
							{{ end }}
						</li>						
					{{ end }}
				</ul>
			</div>
	</body>
</html>`

	return template.Must(template.New("main").Parse(layout))
}

func (p *page) loadMarkdown(path string) {
	md, err := os.ReadFile(path)
	if err != nil {
		p.responseCode = 404
		p.Content = template.HTML("<h1>404 - not found</h1>")
		return
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
