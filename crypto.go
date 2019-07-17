package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Bitcoin block hash.
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

	url := fmt.Sprintf(blockUrl, height-depth)

	data, err = getRead(url)

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

	heading(fmt.Sprintf("Blockchain.Info [depth %d]", depth), url)
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

	url := fmt.Sprintf(blockUrl, stats.Height-depth)

	data, err = getRead(url)

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

	heading(fmt.Sprintf("Moneroblocks.Info [depth %d]", depth), url)
	fmt.Printf("%s\n\n", block.BlockHeader.Hash)

	return nil
}
