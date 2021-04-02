package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

// Target ...
type Target struct {
	Connection string  `json:"connection"`
	Fetch      string  `json:"fetch"`
	Params     []Param `json:"params"`
	Prefetch   bool    `json:"prefetch"`
	Nest       []*Nest `json:"nest"`
	Script     string  `json:"script"`
	Split      *Split  `json:"split"`
	Timezone   string  `json:"timezone"`
	connection *Connection
	extract    string
	prefetch   string
	params     []interface{}
	location   *time.Location
}

// Nest ...
type Nest struct {
	ChildKey   string `json:"childKey"`
	Connection string `json:"connection"`
	Fetch      string `json:"fetch"`
	ParentKey  string `json:"parentKey"`
	Timezone   string `json:"timezone"`
	connection *Connection
	fetch      string
	location   *time.Location
}

// Param ...
type Param struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Default string `json:"default"`
}

// Split ...
type Split struct {
	By     string `json:"by"`
	Layout string `json:"layout"`
	Value  string `json:"value"`
}

func loadTarget() *Target {
	name := path.Join(args.target, defaultTarget)
	checkFileExists(name)

	buffer, err := ioutil.ReadFile(name)
	check(err)

	var target *Target
	json.Unmarshal(buffer, &target)
	resolveTarget(target)
	return target
}

func resolveTarget(target *Target) {
	resolveTargetQuery(target)
	resolveTargetConnection(target)
	resolveTargetLocation(target)
	resolveTargetParams(target)
	resolveTargetNest(target)
	resolveTargetPrefetch(target)
}

func paramName(param Param, i int) string {
	if param.Name != "" {
		return fmt.Sprintf("\"%s\"", param.Name)
	}
	return fmt.Sprintf("#%d", i+1)
}

func applyTargetDefaultParams(target *Target) {
	for i := range target.Params {
		if target.Params[i].Default != "" {
			if len(args.params) < i+1 {
				args.params = append(args.params, target.Params[i].Default)
			}
		}
	}
}

func resolveTargetQuery(target *Target) {
	var err error
	if target.Fetch == "" {
		target.Fetch = defaultExtract
	}
	name := path.Join(args.target, target.Fetch)
	checkFileExists(name)
	target.extract, err = loadFile(name)
	check(err)
}

func resolveTargetConnection(target *Target) {
	if target.Connection == "" {
		target.Connection = defaultConnection
	}
	target.connection = connect(target.Connection)
}

func resolveTargetLocation(target *Target) {
	if target.Timezone != "" {
		var err error
		target.location, err = time.LoadLocation(target.Timezone)
		check(err)
	} else if target.connection.location != nil {
		target.location = target.connection.location
	}
}

func resolveTargetParams(target *Target) {
	applyTargetDefaultParams(target)

	if len(args.params) != len(target.Params) {
		stop("Incorrect number of arguments specified", 1)
	}

	if len(args.params) > 0 {
		target.params = make([]interface{}, len(args.params))
		for i := range args.params {
			target.params[i] = resolveTargetParam(args.params[i], target.Params[i], target.location, i)
		}
	}
}

func resolveTargetParam(text string, param Param, location *time.Location, i int) interface{} {
	var (
		value interface{}
		err   error
	)

	if text == "null" {
		value = nil
	} else if param.Type == "integer" {
		value, err = strconv.Atoi(text)
	} else if param.Type == "date" {
		var t time.Time
		t, err = parseTime(text)
		if err == nil {
			value = formatTimeDb(location, t)
		}
	} else if param.Type == "string" {
		value = text
	} else if param.Type == "boolean" {
		value, err = strconv.ParseBool(text)
	} else if param.Type == "float" {
		value, err = strconv.ParseFloat(text, 32)
	} else {
		stop(fmt.Sprintf("Invalid type \"%s\" specified for param %s", param.Type, paramName(param, i)), 1)
	}

	if err != nil {
		stop(fmt.Sprintf("Unable to convert specified value \"%s\" to type \"%s\" for param %s", text, param.Type, paramName(param, i)), 1)
	}

	return value
}

func resolveTargetNest(target *Target) {
	var err error
	if target.Nest != nil {
		if len(target.Nest) > 0 && target.Nest[0].Fetch == "" {
			target.Nest[0].Fetch = defaultNest
		}

		for i, nest := range target.Nest {
			if nest.Fetch == "" {
				stop(fmt.Sprintf("Unspecified query for nest #%d", i+1), 1)
			}

			name := path.Join(args.target, nest.Fetch)
			checkFileExists(name)
			nest.fetch, err = loadFile(name)
			check(err)

			if nest.Connection != "" {
				nest.connection = connect(nest.Connection)
			} else {
				nest.connection = target.connection
			}

			if nest.Timezone != "" {
				nest.location, err = time.LoadLocation(nest.Timezone)
				check(err)
			} else if target.location != nil {
				nest.location = target.location
			}
		}
	}
}

func resolveTargetPrefetch(target *Target) {
	var err error
	if target.Prefetch {
		name := path.Join(args.target, defaultPrefetch)
		checkFileExists(name)
		target.prefetch, err = loadFile(name)
		check(err)

		if strings.Index(target.extract, "(%s)") == -1 {
			stop("missing \"(%s)\" token in extract query", 1)
		}
	}
}

func validateTarget(target *Target) {
	if target.Split != nil {
		if target.Split.By != "date" {
			stop("Invalid value for \"split\" specified", 1)
		}

		if strings.HasSuffix(args.out, ".json") {
			stop("A file was specified for \"-o\" where a directory was expected", 1)
		}
	}
}
