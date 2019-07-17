package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func main() {
	// date
	fmt.Printf("Date: %s\n\n",
		time.Now().UTC().Format("2006-01-02 15:04 MST"))

	// news feeds
	news()

	// NIST randomness beacons
	nist()

	// BTC block hash
	btc()

	// Monero block hash
	monero()
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title string `xml:"title"`
		} `xml:"item"`
	} `xml:"channel"`
}

func news() {
	var re = regexp.MustCompile(`[^[:ascii:]]+`)

	urls := []string{
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
			log.Fatal(err)
		}

		if len(rss.Channel.Items) < count {
			log.Fatalf("couldn't find %d items", count)
		}

		fmt.Printf("Src: %s (%s)\n ---\n", re.ReplaceAllString(rss.Channel.Title, " "), url)

		for i := 0; i < count; i++ {
			fmt.Printf("%s\n", re.ReplaceAllString(rss.Channel.Items[i].Title, " "))
		}

		fmt.Println()
	}
}

func parseRss(url string) (*Rss, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	rssxml, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	var rss Rss

	if err := xml.Unmarshal(rssxml, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}

func nist() {
	v2URL := "https://beacon.nist.gov/beacon/2.0/pulse/last"

	resp, err := http.Get(v2URL)

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
	btcHeightURL := "https://blockchain.info/q/getblockcount"
	btcBlockURL := "https://blockchain.info/block-height/%d?format=json"

	resp, err := http.Get(btcHeightURL)

	if err != nil {
		log.Fatal(err)
	}

	btcHeight, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	height, err := strconv.ParseInt(string(btcHeight), 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	const depth = 10

	resp, err = http.Get(fmt.Sprintf(btcBlockURL, height-depth))

	if err != nil {
		log.Fatal(err)
	}

	btcBlockJSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}

	var btcBlock struct {
		Blocks []struct {
			Hash string
		}
	}

	if err := json.Unmarshal(btcBlockJSON, &btcBlock); err != nil {
		log.Fatal(err)
	}

	if len(btcBlock.Blocks) == 0 {
		log.Fatal("no blocks found")
	}

	fmt.Printf("Src: Blockchain.Info [block depth %d] (%s)\n ---\n",
		depth, fmt.Sprintf(btcBlockURL, height-depth))

	fmt.Printf("%s\n\n", btcBlock.Blocks[0].Hash)
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
