package stroocts

import "time"

// Structs for selecting data on the main page
type LandingGrouping struct {
	Entries []LandingReport
}

type LandingReport struct {
	Id          string
	Reporter    string
	IssueType   string
	Summary     string
	Full        *string
	Resolved    bool
	Created     *time.Time
	LastUpdated time.Time
}

// Structs for selecting our reporting types
type ReportGrouping struct {
	Types []ReportType
}

type ReportType struct {
	Id   int32
	Text string
}

// Structs for individual reports
type SingleEntry struct {
	PrimaryReport LandingReport
	SubReports    []ChildReports
}

type ChildReports struct {
	Reporter string
	Message  string
	Time     time.Time
}
