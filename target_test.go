package main

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateFile(path string, json string) {
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

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestLoadTargetEmpty(t *testing.T) {
	args = Args{target: "test"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	connections := make(map[string]*Connection)
	connectionName := "default"
	connections[connectionName] = &Connection{
		Driver: "mysql",
		db:     db,
	}
	config = Config{Connections: connections}

	targetPath := path.Join(args.target, defaultTarget)
	CreateFile(targetPath, "{}")
	extractPath := path.Join(args.target, defaultExtract)
	extractConfig := "{}"
	CreateFile(extractPath, extractConfig)

	mock.ExpectBegin()

	target := loadTarget()

	assert.Equal(t, connectionName, target.Connection, "Target Connection set")
	assert.Equal(t, defaultExtract, target.Fetch, "Target Fetch set")
	assert.Equal(t, []Param(nil), target.Params, "Target Params set")
	assert.Equal(t, false, target.Prefetch, "Target Prefetch set")
	assert.Equal(t, []*Nest(nil), target.Nest, "Target Prefetch set")
	assert.Equal(t, "", target.Script, "Target Script set")
	assert.Equal(t, (*Split)(nil), target.Split, "Target Split set")
	assert.Equal(t, "", target.Timezone, "Target Timezone set")
	assert.Equal(t, connections[connectionName], target.connection, "Target connection set")
	assert.Equal(t, extractConfig, target.extract, "Target extract set")
	assert.Equal(t, "", target.prefetch, "Target prefetch set")
	assert.Equal(t, []interface{}(nil), target.params, "Target params set")
	assert.Equal(t, (*time.Location)(nil), target.location, "Target location set")

	DeleteFile(targetPath)
}
