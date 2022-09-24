package main

import (
	"log"
	"net/http"

	"github.com/TeddiO/IncidentProne/src/tmpl"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Landing)
	r.HandleFunc("/new", NewEntry)
	r.HandleFunc("/view/{id}", UpdateEntry)
	r.HandleFunc("/update/{id}", ViewEntry).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func Landing(w http.ResponseWriter, r *http.Request) {
	tmpl.RenderPage("index.html", "templates/index.gohtml", &w)
}

func NewEntry(w http.ResponseWriter, r *http.Request) {
	tmpl.RenderPage("report.html", "templates/report.gohtml", &w)
}

func UpdateEntry(w http.ResponseWriter, r *http.Request) {

}

func ViewEntry(w http.ResponseWriter, r *http.Request) {
	tmpl.RenderPage("viewreport.html", "templates/viewreport.gohtml", &w)
}
