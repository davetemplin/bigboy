package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultInt(priorityInt int, fallbackInt int) int {
	if priorityInt == 0 {
		return fallbackInt
	}
	return priorityInt
}

func getDefaultUint64(priorityInt uint64, fallbackInt uint64) uint64 {
	if priorityInt == 0 {
		return fallbackInt
	}
	return priorityInt
}

func TestParseArgs(t *testing.T) {
	var tests = []struct {
		input  []string
		args   Args
		config Config
	}{
		{[]string{},
			Args{},
			Config{}},

		{[]string{"-n", "-o", "testOut.json", "-q"},
			Args{nulls: true, out: "testOut.json", quiet: true},
			Config{}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.input, " "), func(t *testing.T) {
			args, output := parseArgs("bigboy", tt.input, tt.config)
			assert.Equal(t, "", output, "No unparsed arguments")

			expectedErrors := getDefaultUint64(tt.args.errors, tt.config.Errors)
			expectedPage := getDefaultInt(tt.args.page, tt.config.Page)
			expectedRetries := getDefaultUint64(tt.args.retries, tt.config.Retries)
			expectedWorkers := getDefaultInt(tt.args.workers, tt.config.Workers)

			assert.Equal(t, expectedErrors, args.errors, "Args errors set")
			assert.Equal(t, tt.args.nulls, args.nulls, "Args nulls set")
			assert.Equal(t, tt.args.out, args.out, "Args out set")
			assert.Equal(t, expectedPage, args.page, "Args page set")
			assert.Equal(t, tt.args.params, args.params, "Args params set")
			assert.Equal(t, tt.args.path, args.path, "Args path set")
			assert.Equal(t, tt.args.quiet, args.quiet, "Args quiet set")
			assert.Equal(t, expectedRetries, args.retries, "Args retries set")
			assert.Equal(t, tt.args.target, args.target, "Args target set")
			assert.Equal(t, tt.args.version, args.version, "Args version set")
			assert.Equal(t, expectedWorkers, args.workers, "Args workers set")
		})
	}

}
