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

func TestParseFlags(t *testing.T) {
	var tests = []struct {
		input  []string
		flags  Flags
		config Config
	}{
		{[]string{},
			Flags{},
			Config{}},

		{[]string{},
			Flags{},
			Config{Errors: 555, Nulls: true, Quiet: true, Retries: 11, Workers: 12}},

		{[]string{"-c", "config.json", "-e", "1000", "-n", "-o", "testOut.json", "-p", "100", "-q", "-r", "7", "-v", "-w", "9"},
			Flags{config: "config.json", Args: Args{errors: 1000, nulls: true, out: "testOut.json", page: 100, quiet: true, retries: 7, workers: 9}, version: true},
			Config{}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.input, " "), func(t *testing.T) {
			flags, _ := parseFlags("bigboy", tt.input, &tt.config)

			var expectedConfig string = defaultConfig
			if tt.flags.config != "" {
				expectedConfig = tt.flags.config
			}
			expectedErrors := getDefaultUint64(tt.flags.errors, tt.config.Errors)
			expectedNulls := getDefaultBool(tt.flags.nulls, tt.config.Nulls)
			expectedPage := getDefaultInt(tt.flags.page, tt.config.Page)
			expectedQuiet := getDefaultBool(tt.flags.quiet, tt.config.Quiet)
			expectedRetries := getDefaultUint64(tt.flags.retries, tt.config.Retries)
			expectedWorkers := getDefaultInt(tt.flags.workers, tt.config.Workers)

			assert.Equal(t, expectedConfig, flags.config, "Args config set")
			assert.Equal(t, expectedErrors, flags.errors, "Args errors set")
			assert.Equal(t, expectedNulls, flags.nulls, "Args nulls set")
			assert.Equal(t, tt.flags.out, flags.out, "Args out set")
			assert.Equal(t, expectedPage, flags.page, "Args page set")
			assert.Equal(t, tt.flags.params, flags.params, "Args params set")
			assert.Equal(t, expectedQuiet, flags.quiet, "Args quiet set")
			assert.Equal(t, expectedRetries, flags.retries, "Args retries set")
			assert.Equal(t, tt.flags.target, flags.target, "Args target set")
			assert.Equal(t, tt.flags.version, flags.version, "Args version set")
			assert.Equal(t, expectedWorkers, flags.workers, "Args workers set")
		})
	}

}
