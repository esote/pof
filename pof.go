package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func main() {
	// Date in UTC.
	fmt.Printf("Date: %s\n\n",
		time.Now().UTC().Format("2006-01-02 15:04 MST"))

	sources := [...]func() error{
		news,
		nist,
		btc,
		monero,
	}

	for _, source := range sources {
		if err := source(); err != nil {
			log.Fatal(err)
		}
	}
}

// Structure of an RSS feed, exposing only the fields useful to print news().
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title string `xml:"title"`
		} `xml:"item"`
	} `xml:"channel"`
}

// International news feeds.
func news() error {
	re := regexp.MustCompile(`[^[:ascii:]]+`)

	urls := [...]string{
		"https://www.spiegel.de/international/index.rss",
		"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		"https://feeds.bbci.co.uk/news/world/rss.xml",
		"http://feeds.reuters.com/reuters/worldnews",
		"https://www.economist.com/latest/rss.xml",
	}

	const count = 5

	for _, url := range urls {
		rss, err := parseRss(url)

		if err != nil {
			return err
		}

		if len(rss.Channel.Items) < count {
			return fmt.Errorf("couldn't find %d items", count)
		}

		fmt.Printf("Src: %s (%s)\n ---\n", re.ReplaceAllString(rss.Channel.Title, " "), url)

		for i := 0; i < count; i++ {
			fmt.Printf("%s\n", re.ReplaceAllString(rss.Channel.Items[i].Title, " "))
		}

		fmt.Println()
	}

	return nil
}

// GET and unmarshal specified RSS URL.
func parseRss(url string) (*Rss, error) {
	data, err := getRead(url)

	if err != nil {
		return nil, err
	}

	var rss Rss

	err = xml.Unmarshal(data, &rss)

	return &rss, err
}

// NIST randomness beacon v2.
func nist() error {
	const v2url = "https://beacon.nist.gov/beacon/2.0/pulse/last"

	data, err := getRead(v2url)

	if err != nil {
		return err
	}

	var v2 struct {
		Pulse struct {
			OutputValue string
		}
	}

	if err := json.Unmarshal(data, &v2); err != nil {
		return err
	}

	fmt.Printf("Src: NIST Beacon v2 (%s)\n ---\n", v2url)
	fmt.Printf("%s\n\n", v2.Pulse.OutputValue)

	return nil

}

// BTC block hash.
func btc() error {
	const (
		heightUrl = "https://blockchain.info/q/getblockcount"
		blockUrl  = "https://blockchain.info/block-height/%d?format=json"
		depth     = 10
	)

	data, err := getRead(heightUrl)

	if err != nil {
		return err
	}

	height, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		return err
	}

	data, err = getRead(fmt.Sprintf(blockUrl, height-depth))

	if err != nil {
		return err
	}

	var block struct {
		Blocks []struct {
			Hash string
		}
	}

	if err := json.Unmarshal(data, &block); err != nil {
		return err
	}

	if len(block.Blocks) == 0 {
		return errors.New("no blocks found")
	}

	fmt.Printf("Src: Blockchain.Info [block depth %d] (%s)\n ---\n", depth,
		fmt.Sprintf(blockUrl, height-depth))
	fmt.Printf("%s\n\n", block.Blocks[0].Hash)

	return nil
}

// Monero block hash.
func monero() error {
	const (
		statsUrl = "https://moneroblocks.info/api/get_stats"
		blockUrl = "https://moneroblocks.info/api/get_block_header/%d"
		depth    = 10
	)

	data, err := getRead(statsUrl)

	if err != nil {
		return err
	}

	var stats struct {
		Height int64
	}

	if err := json.Unmarshal(data, &stats); err != nil {
		return err
	}

	if stats.Height < depth {
		return fmt.Errorf("stats.Height < %d", depth)
	}

	data, err = getRead(fmt.Sprintf(blockUrl, stats.Height-depth))

	if err != nil {
		return err
	}

	var block struct {
		BlockHeader struct {
			Hash string
		} `json:"block_header"`
	}

	if err := json.Unmarshal(data, &block); err != nil {
		return err
	}

	fmt.Printf("Src: Moneroblocks.Info [block depth %d] (%s)\n ---\n",
		depth, fmt.Sprintf(blockUrl, stats.Height-depth))
	fmt.Printf("%s\n\n", block.BlockHeader.Hash)

	return nil
}

// Make GET request and read body, reducing duplicate ioutil.ReadAll and error
// checking code.
func getRead(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return data, resp.Body.Close()
}
