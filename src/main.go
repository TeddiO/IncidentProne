package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TeddiO/IncidentProne/src/stroocts"

	"github.com/TeddiO/IncidentProne/src/tmpl"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var (
	// Not the cleanest way of doing this, but for a PoC it's good enough!
	dbConnection *pgx.Conn
)

func init() {
	var err error
	dbConnection, err = pgx.Connect(context.Background(), "postgres://postgres:test@localhost:5432/incidentprone")
	if err != nil {
		panic(err)
	}

	// Use pgxuuid to convert uuids into a string type
	pgxuuid.Register(dbConnection.TypeMap())

}

func main() {
	r := mux.NewRouter()

	// Register some basic routes
	r.HandleFunc("/", Landing)
	r.HandleFunc("/new", NewEntry)
	r.HandleFunc("/create", CreateEntry).Methods("POST")
	r.HandleFunc("/view/{id}", ViewEntry)
	r.HandleFunc("/update/{id}", UpdateEntry).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func Landing(w http.ResponseWriter, r *http.Request) {
	// Query for all of our potential options
	rows, err := dbConnection.Query(context.Background(), "SELECT id::text, \"reporterName\", \"issueType\"::text, \"issueSummary\", resolved, last_updated FROM incidentprone.reports;")
	if err != nil {
		log.Fatal(err)
	}

	var indexData stroocts.LandingGrouping

	for rows.Next() {
		values, err := rows.Values()

		indexData.Entries = append(indexData.Entries, stroocts.LandingReport{Id: values[0].(string), Reporter: values[1].(string), IssueType: values[2].(string),
			Summary: values[3].(string), Resolved: values[4].(bool), LastUpdated: values[5].(time.Time)})

		if err != nil {
			log.Fatal(err)
		}
	}

	// Render our page with all of our entries on them.
	tmpl.RenderPage("index.html", "templates/index.gohtml", &w, indexData)
}

func NewEntry(w http.ResponseWriter, r *http.Request) {

	// Query for all of our potential options
	rows, err := dbConnection.Query(context.Background(), "SELECT * FROM incidentprone.\"reportTypes\";")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our report grouping struct that we plan to dump all of our data in to.
	// We're using a struct here in the event that we plan to maybe send down other data (although we currently don't)
	var dropdownOpts = stroocts.ReportGrouping{}

	// Iterate over our data and cast it to the types we want it to be.
	for rows.Next() {
		values, err := rows.Values()
		dropdownOpts.Types = append(dropdownOpts.Types, stroocts.ReportType{Id: values[1].(int32), Text: values[0].(string)})

		if err != nil {
			log.Fatal(err)
		}
	}

	// Render our page and pass through our array of Types. Normally we'd actually pass the parent struct but
	// as of current we only have our array of Types we we'll just pass that instead.
	tmpl.RenderPage("report.html", "templates/report.gohtml", &w, dropdownOpts.Types)
}

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Need to call this or else we won't populate our fields

	extractedValues := make(map[string]any)

	// Quickly check over our values to ensure they're set and if not - set them
	expectedPostValues := []string{"username", "issueType", "summary", "issue"}
	for _, value := range expectedPostValues {
		if extrValue, isOk := r.PostForm[value]; isOk {
			extractedValues[value] = extrValue[0]
		} else {
			extractedValues[value] = "unset"
		}
	}

	// Special edge case for our resolved check as we're dealing directly with a bool
	alreadyResolvedIssue := false
	if _, isOk := r.PostForm["resolved"]; isOk {
		alreadyResolvedIssue = true
	}

	// Insert our data and fetch back our inserted ID.
	var returnedId string
	if err := dbConnection.QueryRow(context.Background(), "INSERT INTO incidentprone.reports(\"reporterName\", \"issueType\", \"issueSummary\", \"overallIssue\", resolved) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		extractedValues["username"], extractedValues["issueType"], extractedValues["summary"], extractedValues["issue"], alreadyResolvedIssue).Scan(&returnedId); err != nil {
		log.Fatal(err)
	}

	// And then redirect our user to their report.
	http.Redirect(w, r, fmt.Sprintf("/view/%s", returnedId), http.StatusPermanentRedirect)
}

func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
}

func ViewEntry(w http.ResponseWriter, r *http.Request) {
	var reportEntry stroocts.SingleEntry
	vars := mux.Vars(r)

	var initialReport stroocts.LandingReport

	if err := dbConnection.QueryRow(context.Background(), "SELECT id::text, \"reporterName\", \"issueType\"::text, \"issueSummary\", \"overallIssue\", resolved, last_updated, created FROM incidentprone.reports WHERE id = $1;", vars["id"]).Scan(
		&initialReport.Id, &initialReport.Reporter, &initialReport.IssueType, &initialReport.Summary, &initialReport.Full, &initialReport.Resolved, &initialReport.LastUpdated, &initialReport.Created); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	reportEntry.PrimaryReport = initialReport

	fmt.Println(initialReport)

	tmpl.RenderPage("viewreport.html", "templates/viewreport.gohtml", &w, reportEntry)
}
