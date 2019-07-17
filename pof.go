package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func main() {
	// Date in UTC.
	fmt.Printf("Date: %s\n\n",
		time.Now().UTC().Format("2006-01-02 15:04 MST"))

	sources := [...]func() error{
		arxiv,
		btc,
		monero,
		news,
		nist,
	}

	for _, source := range sources {
		if err := source(); err != nil {
			log.Fatal(err)
		}
	}
}

var nonascii = regexp.MustCompile(`[^[:ascii:]]+`)

func heading(title string, url string) {
	fmt.Printf("Src: %s (%s)\n ---\n", title, url)
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
