package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// International news feeds.
func news() error {
	const count = 5

	type source struct {
		name string
		url  string
	}

	sources := [...]source{
		{
			"Der Spiegel international",
			"https://www.spiegel.de/international/index.rss",
		},
		{
			"New York Times world news",
			"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		},
		{
			"BBC world news",
			"https://feeds.bbci.co.uk/news/world/rss.xml",
		},
		{
			"The Economist latest updates",
			"https://www.economist.com/latest/rss.xml",
		},
	}

	for _, source := range sources {
		data, err := getRead(source.url)

		if err != nil {
			return err
		}

		// Structure of RSS feed, exposing only the useful fields
		var rss struct {
			Name   string   `xml:"channel>title"`
			Titles []string `xml:"channel>item>title"`
		}

		if err = xml.Unmarshal(data, &rss); err != nil {
			return err
		}

		if len(rss.Titles) < count {
			return fmt.Errorf("couldn't find %d articles", count)
		}

		heading(source.name, source.url)

		for i := 0; i < count; i++ {
			fmt.Println(strings.TrimSpace(nonascii.ReplaceAllString(rss.Titles[i], " ")))
		}

		fmt.Println()
	}

	return nil
}
