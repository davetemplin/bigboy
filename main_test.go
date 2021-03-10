package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultBool(priorityBool bool, fallbackBool bool) bool {
	if !priorityBool {
		return fallbackBool
	}
	return priorityBool
}

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

		{[]string{},
			Args{},
			Config{Errors: 555, Nulls: true, Quiet: true, Retries: 11, Workers: 12}},

		{[]string{"-c", "config.json", "-e", "1000", "-n", "-o", "testOut.json", "-p", "100", "-q", "-r", "7", "-v", "-w", "9"},
			Args{config: "config.json", errors: 1000, nulls: true, out: "testOut.json", page: 100, quiet: true, retries: 7, version: true, workers: 9},
			Config{}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.input, " "), func(t *testing.T) {
			args := *parseArgs("bigboy", tt.input, tt.config)

			var expectedConfig string = defaultConfig
			if tt.args.config != "" {
				expectedConfig = tt.args.config
			}
			expectedErrors := getDefaultUint64(tt.args.errors, tt.config.Errors)
			expectedNulls := getDefaultBool(tt.args.nulls, tt.config.Nulls)
			expectedPage := getDefaultInt(tt.args.page, tt.config.Page)
			expectedQuiet := getDefaultBool(tt.args.quiet, tt.config.Quiet)
			expectedRetries := getDefaultUint64(tt.args.retries, tt.config.Retries)
			expectedWorkers := getDefaultInt(tt.args.workers, tt.config.Workers)

			assert.Equal(t, expectedConfig, args.config, "Args config set")
			assert.Equal(t, expectedErrors, args.errors, "Args errors set")
			assert.Equal(t, expectedNulls, args.nulls, "Args nulls set")
			assert.Equal(t, tt.args.out, args.out, "Args out set")
			assert.Equal(t, expectedPage, args.page, "Args page set")
			assert.Equal(t, tt.args.params, args.params, "Args params set")
			assert.Equal(t, expectedQuiet, args.quiet, "Args quiet set")
			assert.Equal(t, expectedRetries, args.retries, "Args retries set")
			assert.Equal(t, tt.args.target, args.target, "Args target set")
			assert.Equal(t, tt.args.version, args.version, "Args version set")
			assert.Equal(t, expectedWorkers, args.workers, "Args workers set")
		})
	}

}
