package main

import (
	"encoding/json"
	"fmt"
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

func getDefaultConfig() Config {
	return Config{
		Errors:  defaultErrors,
		Page:    defaultPage,
		Retries: defaultRetries,
		Workers: defaultWorkers,
	}
}

func loadConfig(path string) Config {
	var c Config = getDefaultConfig()
	if fileExists(path) {
		buffer, err := ioutil.ReadFile(path)
		check(err)
		json.Unmarshal(buffer, &c)
	} else {
		fmt.Println("Using default configuration")
	}
	return c
}
