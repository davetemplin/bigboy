package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func TestRun(t *testing.T) {
	args = Args{out: "test_main.json", page: 10, workers: 1, target: "test"}
	os.Remove(args.out)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	conn := &Connection{
		Driver:   "mysql",
		Server:   "tcp(mysql-test.ac.uk:4496)",
		Database: "test",
		// Dsn: "",
		Port:     4496,
		User:     "TESTER",
		Password: "PASS",
		Timezone: "UTC",
		Max:      1,
		db:       db,
		// location: *time.Location,
	}

	target := &Target{
		Connection: "conn",
		Fetch:      "",
		Params:     []Param(nil),
		Prefetch:   false,
		Nest:       []*Nest(nil),
		Script:     "",
		Split:      (*Split)(nil),
		Timezone:   "UTC",
		connection: conn,
		extract:    "SELECT name, date FROM testTable",
		prefetch:   "",
		params:     []interface{}(nil),
		location:   (*time.Location)(nil),
	}

	type TestQuery struct {
		One string `json:"One"`
		Two int32  `json:"Two"`
	}
	var expectedOutput = TestQuery{
		One: "ValueOne",
		Two: 2,
	}
	var testOut TestQuery

	mock.ExpectQuery(target.extract).WillReturnRows(mock.NewRows([]string{"One", "Two"}).AddRow(expectedOutput.One, expectedOutput.Two))

	run(target)
	disconnect()

	data, err := ioutil.ReadFile(args.out)
	assert.Equal(t, nil, err, "Can read from output file")
	err = json.Unmarshal(data, &testOut)
	assert.Equal(t, nil, err, "Can read output json")
	assert.Equal(t, expectedOutput, testOut, "Output matches input")
	os.Remove(args.out)
}

func TestRunPrefetch(t *testing.T) {
	args = Args{out: "test_main.json", page: 10, workers: 1, target: "test", errors: 3}
	os.Remove(args.out)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	conn := &Connection{
		Driver:   "mysql",
		Server:   "tcp(mysql-test.ac.uk:4496)",
		Database: "test",
		// Dsn: "",
		Port:     4496,
		User:     "TESTER",
		Password: "PASS",
		Timezone: "UTC",
		Max:      1,
		db:       db,
		// location: *time.Location,
	}

	target := &Target{
		Connection: "conn",
		Fetch:      "",
		Params:     []Param(nil),
		Prefetch:   true,
		Nest:       []*Nest(nil),
		Script:     "",
		Split:      (*Split)(nil),
		Timezone:   "UTC",
		connection: conn,
		extract:    "SELECT id, name FROM testTable WHERE id IN (%s)",
		prefetch:   "SELECT id FROM testTable WHERE date >= \"2021-03-24\"",
		params:     []interface{}(nil),
		location:   (*time.Location)(nil),
	}

	type TestQuery struct {
		One string `json:"One"`
		Two int32  `json:"Two"`
	}
	expectedOutput := []TestQuery{{One: "row1", Two: 12}, {One: "row2", Two: 34}}
	prefetchOutput := []uint64{12, 34}
	var testOut []TestQuery

	expectedQuery := fmt.Sprintf(target.extract, "?")
	mock.ExpectQuery(target.prefetch).WillReturnRows(mock.NewRows([]string{"Id"}).AddRow(prefetchOutput[0]).AddRow(prefetchOutput[1]))
	mock.ExpectQuery(expectedQuery).WillReturnRows(mock.NewRows([]string{"One", "Two"}).AddRow(expectedOutput[0].One, expectedOutput[0].Two).AddRow(expectedOutput[1].One, expectedOutput[1].Two))

	run(target)
	disconnect()

	data, err := ioutil.ReadFile(args.out)
	assert.Equal(t, nil, err, "Can read from output file")

	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var testQuery TestQuery
		err = decoder.Decode(&testQuery)
		assert.Equal(t, nil, err, "Can read output json")
		testOut = append(testOut, testQuery)
	}

	assert.Equal(t, expectedOutput, testOut, "Output matches input")
	os.Remove(args.out)
}
