package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	Connections map[string]*Connection `json:"connections"`
	Errors      uint64                 `json:"errors"`
	Nulls       bool                   `json:"nulls"`
	Page        int                    `json:"page"`
	Quiet       bool                   `json:"quiet"`
	Retries     uint64                 `json:"retries"`
	Workers     int                    `json:"workers"`
}

var config Config

const undefined = ^uint64(0)

func loadConfig(path string) {
	config.Errors = undefined
	config.Retries = undefined

	if fileExists(path) {
		buffer, err := ioutil.ReadFile(path)
		check(err)
		json.Unmarshal(buffer, &config)
	}

	if config.Errors == undefined {
		config.Errors = 100
	}
	if config.Page == 0 {
		config.Page = 1000
	}
	if config.Retries == undefined {
		config.Retries = 3
	}
	if config.Workers == 0 {
		config.Workers = 4
	}
}
