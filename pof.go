package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	parser := gofeed.NewParser()

	// date
	fmt.Printf("Date: %s\n\n",
		time.Now().UTC().Format("2006-01-02 15:04 MST"))

	// news feeds
	news(parser, 5)

	// NIST randomness beacons
	nist()

	// BTC block hash
	btc()

	// Monero block hash
	monero()
}

func news(parser *gofeed.Parser, count int) {
	urls := []string{
		"https://www.spiegel.de/international/index.rss",
		"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		"https://feeds.bbci.co.uk/news/world/rss.xml",
		"http://feeds.reuters.com/reuters/worldnews",
		"https://www.economist.com/international/rss.xml",
	}

	for _, url := range urls {
		feed, err := parser.ParseURL(url)

		if err != nil {
			log.Fatal(err)
		}

		if len(feed.Items) < count {
			log.Fatalf("couldn't find %d items", count)
		}

		fmt.Printf("Src: %s (%s)\n ---\n", feed.Title, url)

		for i := 0; i < count; i++ {
			fmt.Printf("%s\n", feed.Items[i].Title)
		}

		fmt.Println()
	}
}

func nist() {
	v1URL := "https://beacon.nist.gov/rest/record/last"
	v2URL := "https://beacon.nist.gov/beacon/2.0/pulse/last"

	resp, err := http.Get(v1URL)

	if err != nil {
		log.Fatal(err)
	}

	v1XML, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var v1 struct {
		XMLName     xml.Name `xml:"record"`
		OutputValue string   `xml:"outputValue"`
	}

	if err := xml.Unmarshal(v1XML, &v1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Src: NIST Beacon v1 (%s)\n ---\n", v1URL)
	fmt.Printf("%s\n\n", v1.OutputValue)

	resp, err = http.Get(v2URL)

	if err != nil {
		log.Fatal(err)
	}

	v2JSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var v2 struct {
		Pulse struct {
			OutputValue string
		}
	}

	if err := json.Unmarshal(v2JSON, &v2); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Src: NIST Beacon v2 (%s)\n ---\n", v2URL)
	fmt.Printf("%s\n\n", v2.Pulse.OutputValue)

}

func btc() {
	btcURL := "https://blockchain.info/blocks/?format=json"

	resp, err := http.Get(btcURL)

	if err != nil {
		log.Fatal(err)
	}

	btcJSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var btc struct {
		Blocks []struct {
			Hash string
		}
	}

	if err := json.Unmarshal(btcJSON, &btc); err != nil {
		log.Fatal(err)
	}

	depth := 10

	if len(btc.Blocks) < depth {
		log.Fatalf("len(btc.Blocks) < %d", depth)
	}

	fmt.Printf("Src: Blockchain.Info [block depth %d] (%s)\n ---\n",
		depth, btcURL)
	fmt.Printf("%s\n\n", btc.Blocks[10].Hash)
}

func monero() {
	monStatURL := "https://moneroblocks.info/api/get_stats"
	monURL := "https://moneroblocks.info/api/get_block_header/%d"

	resp, err := http.Get(monStatURL)

	if err != nil {
		log.Fatal(err)
	}

	monJSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var monStats struct {
		Height int64
	}

	if err := json.Unmarshal(monJSON, &monStats); err != nil {
		log.Fatal(err)
	}

	depth := int64(10)

	if monStats.Height < depth {
		log.Fatalf("monStats.Height < %d", depth)
	}

	resp, err = http.Get(fmt.Sprintf(monURL, monStats.Height-depth))

	if err != nil {
		log.Fatal(err)
	}

	monJSON, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var monBlockHdr struct {
		BlockHeader struct {
			Hash string
		} `json:"block_header"`
	}

	if err := json.Unmarshal(monJSON, &monBlockHdr); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Src: Moneroblocks.Info [block depth %d] (%s)\n ---\n",
		depth,
		fmt.Sprintf(monURL, monStats.Height-depth))
	fmt.Printf("%s\n\n", monBlockHdr.BlockHeader.Hash)

}
