package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLoadConfigMissing(t *testing.T) {
	configPath := "test_config.json"
	config := loadConfig(configPath)

	assert.Equal(t, map[string]*Connection(map[string]*Connection(nil)), config.Connections, "Config Connections set")
	assert.Equal(t, toUint64(defaultErrors), config.Errors, "Config Errors set")
	assert.Equal(t, false, config.Nulls, "Config Nulls set")
	assert.Equal(t, defaultPage, config.Page, "Config Page set")
	assert.Equal(t, false, config.Quiet, "Config Quiet set")
	assert.Equal(t, toUint64(defaultRetries), config.Retries, "Config Retries set")
	assert.Equal(t, defaultWorkers, config.Workers, "Config Workers set")
}

func TestLoadConfigEmpty(t *testing.T) {
	configPath := "test_config.json"
	CreateConfig(configPath, "{}")

	config := loadConfig(configPath)

	assert.Equal(t, map[string]*Connection(map[string]*Connection(nil)), config.Connections, "Config Connections set")
	assert.Equal(t, toUint64(defaultErrors), config.Errors, "Config Errors set")
	assert.Equal(t, false, config.Nulls, "Config Nulls set")
	assert.Equal(t, defaultPage, config.Page, "Config Page set")
	assert.Equal(t, false, config.Quiet, "Config Quiet set")
	assert.Equal(t, toUint64(defaultRetries), config.Retries, "Config Retries set")
	assert.Equal(t, defaultWorkers, config.Workers, "Config Workers set")

	DeleteConfig(configPath)
}

func TestLoadConfigOverride(t *testing.T) {
	configPath := "test_config.json"
	testConfig := &Config{
		Errors:  toUint64(10),
		Nulls:   true,
		Page:    10,
		Quiet:   true,
		Retries: toUint64(10),
		Workers: 10,
	}
	json, err := json.Marshal(testConfig)
	assert.Equal(t, err, nil, "JSON config created")

	CreateConfig(configPath, string(json))

	config := loadConfig(configPath)

	assert.Equal(t, map[string]*Connection(map[string]*Connection(nil)), config.Connections, "Config Connections set")
	assert.Equal(t, (*testConfig).Errors, config.Errors, "Config Errors set")
	assert.Equal(t, (*testConfig).Nulls, config.Nulls, "Config Nulls set")
	assert.Equal(t, (*testConfig).Page, config.Page, "Config Page set")
	assert.Equal(t, (*testConfig).Quiet, config.Quiet, "Config Quiet set")
	assert.Equal(t, (*testConfig).Retries, config.Retries, "Config Retries set")
	assert.Equal(t, (*testConfig).Workers, config.Workers, "Config Workers set")

	DeleteConfig(configPath)
}
