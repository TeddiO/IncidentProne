package tmpl

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

var (
	// Embed our templates into the app
	//go:embed templates/*
	embededFiles embed.FS
	templates    map[string]*template.Template
)

func init() {

	// Make sure we're initialized...
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	// Read the directory to make sure it exists
	templFiles, err := fs.ReadDir(embededFiles, "templates")
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over and ensure we're not dealing with any subdirs
	for _, templ := range templFiles {
		if templ.IsDir() {
			continue
		}

		// Parse the template from the virtual file system
		parsedTemplate, err := template.ParseFS(embededFiles, fmt.Sprintf("%s/%s", "templates", templ.Name()))
		if err != nil {
			log.Fatal(err)
		}

		// And then queue it up in our map so we can easily reference it later.
		templates[templ.Name()] = parsedTemplate
	}
}

func RenderPage(name string, resp *http.ResponseWriter, data any) {
	targetTemplate, exists := templates[name]
	if !exists {
		log.Fatalf("Attempted to specify %s template which doesn't exist\n", name)
	}

	if err := targetTemplate.Execute(*resp, data); err != nil {
		fmt.Println(err)
	}
}
