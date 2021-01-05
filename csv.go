package main

import (
	"encoding/csv"
	"log"
	"os"
	"io"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type Issue struct {
	Assignee string `json:"assignee,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Body string `json:"body,omitempty"`
	Closed bool `json:"closed"`
	Labels []int `json:"labels,omitempty"`
	Due_date string `json:"due_date,omitempty"`
	Title string `json:"title"`
}

var (
	accessToken = ""
	apiUrl = "https://git.freifunk-franken.de/api"
	repo = "freifunk-franken/firmware"
	mantisUrl = "https://mantis.freifunk-franken.de"
)

// from: https://gist.github.com/drernie/5684f9def5bee832ebc50cabb46c377a
func csvToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

func readCsv() []Issue {
	var issues []Issue
	_ = issues

	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	csv := csvToMap(file)

	for _, s := range(csv) {
		var issue Issue
		issue.Title = s["Summary"]
		reporter := "Reported by: " + s["Reporter"]
		imported := "Imported from: " + mantisUrl + "/view.php?id=" + s["Id"]
		date := "Submitted on: " + s["Date Submitted"]
		issue.Body = s["Description"] + "\n" + "\n" + reporter + "\n" + date + "\n" + imported + "\n"
		issue.Closed = false
		issue.Labels = []int{ 28 }
		if (s["Assigned To"] != "") {
			assgn := s["Assigned To"]
			if (assgn == "Adrian Schmutzler") { assgn = "adschm" }
			issue.Assignees = []string{ assgn }
		}

		issues = append(issues, issue)
	}

	return issues
}

func sendIssues(issues []Issue) {
	for _, s := range(issues) {
		jsonResult, err := json.Marshal(s)
		_ = jsonResult
		if err != nil {
			log.Print(err)
			continue
		}
		resp, err := http.Post(apiUrl + "/v1/repos/" + repo + "/issues?access_token=" +  accessToken, "application/json", bytes.NewReader(jsonResult))
		if err != nil {
			log.Print(err)
		}
		log.Println(resp)
		fmt.Println(string(jsonResult))
	}
}

func main() {
	/*var issues Issue[];
	issues = getIssues();*/
	issues := readCsv()
	sendIssues(issues)
}
