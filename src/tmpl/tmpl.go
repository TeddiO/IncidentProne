package tmpl

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderPage(name string, path string, resp *http.ResponseWriter, data any) {
	// Wouldn't usually handle templates like this, but for on the fly refreshing it's useful!
	template.New(name)
	returnedTemplate, _ := template.ParseFiles(path)

	if err := returnedTemplate.Execute(*resp, data); err != nil {
		fmt.Println(err)
	}
}
