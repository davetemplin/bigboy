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

func TestLoadTargetEmpty(t *testing.T) {
	args = Args{target: "testEmpty"}
	targetPath := path.Join(args.target, defaultTarget)
	extractPath := path.Join(args.target, defaultExtract)

	os.Remove(targetPath)
	os.Remove(extractPath)
	os.Remove(args.target)

	err := os.Mkdir(args.target, 0755)
	if err != nil {
		t.Fatal(err)
	}

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

	CreateFile(targetPath, "{}")
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

	os.Remove(targetPath)
	os.Remove(extractPath)
	os.Remove(args.target)
}

func TestLoadTarget(t *testing.T) {
	args = Args{target: "test"}
	targetPath := path.Join(args.target, defaultTarget)
	extractPath := path.Join(args.target, defaultExtract)
	prefetchPath := path.Join(args.target, defaultPrefetch)

	os.Remove(targetPath)
	os.Remove(extractPath)
	os.Remove(prefetchPath)
	os.Remove(args.target)

	err := os.Mkdir(args.target, 0755)
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	connections := make(map[string]*Connection)
	connectionName := "conn"
	connections[connectionName] = &Connection{
		Driver:   "mysql",
		Server:   "tcp(mysql-test.ac.uk:4496)",
		Database: "test",
		// Dsn: "",
		Port:     4496,
		User:     "TESTER",
		Password: "PASS",
		Timezone: "UTC",
		Max:      3,
		db:       db,
		// location: *time.Location,
	}
	config = Config{Connections: connections}

	CreateFile(targetPath, "{\"connection\":\"conn\",\"prefetch\":true}")
	extractConfig := "SELECT name, date FROM testTable WHERE id IN (%s)"
	CreateFile(extractPath, extractConfig)
	prefetchConfig := "SELECT id FROM testTable"
	CreateFile(prefetchPath, prefetchConfig)

	mock.ExpectBegin()

	target := loadTarget()

	assert.Equal(t, connectionName, target.Connection, "Target Connection set")
	assert.Equal(t, defaultExtract, target.Fetch, "Target Fetch set")
	assert.Equal(t, []Param(nil), target.Params, "Target Params set")
	assert.Equal(t, true, target.Prefetch, "Target Prefetch set")
	assert.Equal(t, []*Nest(nil), target.Nest, "Target Nest set")
	assert.Equal(t, "", target.Script, "Target Script set")
	assert.Equal(t, (*Split)(nil), target.Split, "Target Split set")
	assert.Equal(t, "", target.Timezone, "Target Timezone set")
	assert.Equal(t, connections[connectionName], target.connection, "Target connection set")
	assert.Equal(t, extractConfig, target.extract, "Target extract set")
	assert.Equal(t, prefetchConfig, target.prefetch, "Target prefetch set")
	assert.Equal(t, []interface{}(nil), target.params, "Target params set")
	assert.Equal(t, (*time.Location)(nil), target.location, "Target location set")

	os.Remove(targetPath)
	os.Remove(extractPath)
	os.Remove(prefetchPath)
	os.Remove(args.target)
}
