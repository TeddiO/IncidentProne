package tmpl

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderPage(name string, path string, resp *http.ResponseWriter) {
	template.New(name)
	returnedTemplate, _ := template.ParseFiles(path)

	if err := returnedTemplate.Execute(*resp, nil); err != nil {
		fmt.Println(err)
	}
}

// func RenderPage(name string, template template.Template, resp *http.ResponseWriter) {
// 	if err := template.Execute(*resp, nil); err != nil {
// 		fmt.Println(err)
// 	}
// }
