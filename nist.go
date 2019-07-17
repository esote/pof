package main

import (
	"encoding/json"
	"fmt"
)

// NIST beacon v2.
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

	heading("NIST beacon v2", v2url)
	fmt.Printf("%s\n\n", v2.Pulse.OutputValue)

	return nil

}
