package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/TeddiO/IncidentProne/src/tmpl"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var (
	dbConnection *pgx.Conn
)

func init() {
	var err error
	dbConnection, err = pgx.Connect(context.Background(), "postgres://postgres:test@localhost:5432/incidentprone")
	if err != nil {
		panic(err)
	}
}

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
	tmpl.RenderPage("index.html", "templates/index.gohtml", &w, nil)
}

type reportGrouping struct {
	Types []reportType
}

type reportType struct {
	Id   int32
	Text string
}

func NewEntry(w http.ResponseWriter, r *http.Request) {

	// Query for all of our potential options
	rows, err := dbConnection.Query(context.Background(), "SELECT * FROM incidentprone.\"reportTypes\";")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our report grouping struct that we plan to dump all of our data in to.
	// We're using a struct here in the event that we plan to maybe send down other data (although we currently don't)
	var dropdownOpts = reportGrouping{}

	// Iterate over our data and cast it to the types we want it to be.
	for rows.Next() {
		values, err := rows.Values()
		dropdownOpts.Types = append(dropdownOpts.Types, reportType{Id: values[1].(int32), Text: values[0].(string)})

		if err != nil {
			log.Fatal(err)
		}
	}

	// Render our page and pass through our array of Types. Normally we'd actually pass the parent struct but
	// as of current we only have our array of Types we we'll just pass that instead.
	tmpl.RenderPage("report.html", "templates/report.gohtml", &w, dropdownOpts.Types)
}

func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)

}

func ViewEntry(w http.ResponseWriter, r *http.Request) {
	tmpl.RenderPage("viewreport.html", "templates/viewreport.gohtml", &w, nil)
}
