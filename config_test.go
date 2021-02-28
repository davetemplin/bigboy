package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func CreateConfig(path string, json string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(json)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func DeleteConfig(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestLoadConfigEmpty(t *testing.T) {
	configPath := "test_config.json"
	CreateConfig(configPath, "{}")

	loadConfig(configPath)

	assert.Equal(t, config.Connections, map[string]*Connection(map[string]*Connection(nil)), "Config has Connections not set")
	assert.Equal(t, config.Errors, toUint64(100), "Config has Errors set")
	assert.Equal(t, config.Nulls, false, "Config has Nulls not set")
	assert.Equal(t, config.Page, 1000, "Config has Page set")
	assert.Equal(t, config.Quiet, false, "Config has Quiet not set")
	assert.Equal(t, config.Retries, toUint64(3), "Config has Retries set")
	assert.Equal(t, config.Workers, 4, "Config has Workers set")

	DeleteConfig(configPath)
}
