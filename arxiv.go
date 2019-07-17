package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// arXiv recently submitted preprints.
func arxiv() error {
	const (
		queryUrl = "https://export.arxiv.org/api/query?" +
			"search_query=all&sortBy=submittedDate&" +
			"sortOrder=descending&max_results=%d"
		count  = 10
		maxlen = 80
	)

	url := fmt.Sprintf(queryUrl, count)

	data, err := getRead(url)

	if err != nil {
		return err
	}

	var arxiv struct {
		Entries []struct {
			Published string   `xml:"published"`
			Title     string   `xml:"title"`
			Authors   []string `xml:"author>name"`
		} `xml:"entry"`
	}

	if err := xml.Unmarshal(data, &arxiv); err != nil {
		return err
	}

	if len(arxiv.Entries) != count {
		return fmt.Errorf("response length mismatched %d", count)
	}

	heading("arXiv", url)

	for _, entry := range arxiv.Entries {
		entry.Title = nonascii.ReplaceAllString(entry.Title, " ")
		entry.Title = strings.ReplaceAll(entry.Title, "\n ", "")

		if len(entry.Title) > maxlen {
			entry.Title = entry.Title[:maxlen] + "..."
		}

		entry.Title = strings.TrimSpace(entry.Title)

		var author string

		if len(entry.Authors) > 0 {
			author = entry.Authors[0]
		}

		if len(entry.Authors) > 1 {
			author += ", et al."
		}

		author = strings.TrimSpace(nonascii.ReplaceAllString(author, " "))

		fmt.Printf("%s (%s, %s)\n", entry.Title, author,
			entry.Published)
	}

	fmt.Println()

	return nil
}
