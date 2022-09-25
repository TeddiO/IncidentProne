package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TeddiO/IncidentProne/src/stroocts"

	"github.com/TeddiO/IncidentProne/src/tmpl"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var (
	// Not the cleanest way of doing this, but for a PoC it's good enough!
	dbConnection *pgx.Conn
)

func init() {
	fmt.Println("Booting initial database settings...")
	dbHostname := os.Getenv("DB_HOSTNAME")

	var err error
	dbConnection, err = pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:test@%s:5432/incidentprone", dbHostname))
	if err != nil {
		panic(err)
	}
	fmt.Println("Booted all database settings!")

}

func main() {
	fmt.Println("registering router...")
	r := mux.NewRouter()

	// Register some basic routes
	fmt.Println("Creating routes...")
	r.HandleFunc("/", Landing)
	r.HandleFunc("/new", NewEntry)
	r.HandleFunc("/create", CreateEntry).Methods("POST")
	r.HandleFunc("/view/{id}", ViewEntry)
	r.HandleFunc("/update", UpdateEntry).Methods("POST")

	fmt.Println("Booting server")
	broadcastHost := os.Getenv("APP_HOSTNAME")
	broadcastPort := os.Getenv("APP_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", broadcastHost, broadcastPort), r))
}

func Landing(w http.ResponseWriter, r *http.Request) {
	// Query for all of our potential options
	rows, err := dbConnection.Query(context.Background(), "SELECT id::text, \"reporterName\", \"reason\", \"issueSummary\", resolved, last_updated FROM incidentprone.reports JOIN incidentprone.\"reportTypes\" ON \"issueType\" = \"internalId\";")
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
	tmpl.RenderPage("index.gohtml", &w, indexData)
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
	tmpl.RenderPage("report.gohtml", &w, dropdownOpts.Types)
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
	r.ParseForm()

	// Insert our sub message wqithout any real care for errors :')
	dbConnection.Exec(context.Background(), "INSERT INTO incidentprone.sub_reports(username, message, referenced_issue) VALUES ($1, $2, $3)",
		r.PostForm["username"][0], r.PostForm["issue"][0], r.PostForm["issueId"][0])

	// And update the last updated field to ensure we can see when something was last updated!
	dbConnection.Exec(context.Background(), "UPDATE incidentprone.reports SET last_updated = CURRENT_TIMESTAMP WHERE id = $1;", r.PostForm["issueId"][0])

	// And similar thing with if we're marking it as resolved.
	if _, isOk := r.PostForm["resolved"]; isOk {
		dbConnection.Exec(context.Background(), "UPDATE incidentprone.reports SET resolved = True WHERE id = $1;", r.PostForm["issueId"][0])
	}

	// Then to reload the page to show change, send them back to the page.
	http.Redirect(w, r, fmt.Sprintf("/view/%s", r.PostForm["issueId"][0]), http.StatusTemporaryRedirect)
}

func ViewEntry(w http.ResponseWriter, r *http.Request) {
	var reportEntry stroocts.SingleEntry
	vars := mux.Vars(r)
	var initialReport stroocts.LandingReport

	// Grab our main report data. If we have an error for whatever reason, safely back out to our index
	if err := dbConnection.QueryRow(context.Background(), "SELECT id::text, \"reporterName\", reason, \"issueSummary\", \"overallIssue\", resolved, last_updated, created FROM incidentprone.reports JOIN incidentprone.\"reportTypes\" ON \"issueType\" = \"internalId\" WHERE id = $1;", vars["id"]).Scan(
		&initialReport.Id, &initialReport.Reporter, &initialReport.IssueType, &initialReport.Summary, &initialReport.Full, &initialReport.Resolved, &initialReport.LastUpdated, &initialReport.Created); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	if initialReport.Resolved {
		timeToResolve := initialReport.LastUpdated.Sub(*initialReport.Created).Round(time.Second)
		initialReport.TotalTime = &timeToResolve
	}

	reportEntry.PrimaryReport = initialReport

	// Then load up our sub messages said the incident.
	rows, err := dbConnection.Query(context.Background(), "SELECT username, message, time FROM incidentprone.sub_reports WHERE referenced_issue = $1;", vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Loop-de-loop out the values and cast them.
	for rows.Next() {
		values, err := rows.Values()
		reportEntry.SubReports = append(reportEntry.SubReports, stroocts.ChildReports{Reporter: values[0].(string), Message: values[1].(string), Time: values[2].(time.Time)})

		if err != nil {
			log.Fatal(err)
		}
	}

	// And render us out!
	tmpl.RenderPage("viewreport.gohtml", &w, reportEntry)
}
