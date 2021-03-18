package main

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateTarget(path string, json string) {
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

func DeleteTarget(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestLoadTarget(t *testing.T) {
	targetPath := "test"
	filePath := path.Join(targetPath, defaultTarget)
	CreateTarget(filePath, "{}")

	target := loadTarget(targetPath)

	assert.Equal(t, "", target.Connection, "Target Connection set")
	assert.Equal(t, "", target.Fetch, "Target Fetch set")
	assert.Equal(t, []Param(nil), target.Params, "Target Params set")
	assert.Equal(t, false, target.Prefetch, "Target Prefetch set")
	assert.Equal(t, []*Nest(nil), target.Nest, "Target Prefetch set")
	assert.Equal(t, "", target.Script, "Target Script set")
	assert.Equal(t, (*Split)(nil), target.Split, "Target Split set")
	assert.Equal(t, "", target.Timezone, "Target Timezone set")
	assert.Equal(t, (*Connection)(nil), target.connection, "Target connection set")
	assert.Equal(t, "", target.extract, "Target extract set")
	assert.Equal(t, "", target.prefetch, "Target prefetch set")
	assert.Equal(t, []interface{}(nil), target.params, "Target params set")
	assert.Equal(t, (*time.Location)(nil), target.location, "Target location set")

	DeleteTarget(filePath)
}
