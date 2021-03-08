package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	var tests = []struct {
		input []string
		args  Args
	}{
		{[]string{"-n", "-o", "testOut.json", "-q"},
			Args{nulls: true, out: "testOut.json", quiet: true}},

		// {[]string{"-n", "true"},
		//		Args{nulls: true}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.input, " "), func(t *testing.T) {
			args, output := parseArgs("bigboy", tt.input)
			assert.Equal(t, "", output, "No unparsed arguments")
			// if !reflect.DeepEqual(*args, tt.args) {
			// 	t.Errorf("conf got %+v, want %+v", *args, tt.args)
			// }
			// assert.Equal(t, toUint64(defaultErrors), args.errors, "Args errors set")
			assert.Equal(t, tt.args.nulls, args.nulls, "Args nulls set")
			assert.Equal(t, tt.args.out, args.out, "Args out set")
			// assert.Equal(t, defaultPage, args.page, "Args page set")
			assert.Equal(t, tt.args.params, args.params, "Args params set")
			assert.Equal(t, tt.args.path, args.path, "Args path set")
			assert.Equal(t, tt.args.quiet, args.quiet, "Args quiet set")
			// assert.Equal(t, toUint64(defaultRetries), args.retries, "Args retries set")
			assert.Equal(t, tt.args.target, args.target, "Args target set")
			assert.Equal(t, tt.args.version, args.version, "Args version set")
			// assert.Equal(t, defaultWorkers, args.workers, "Args workers set")
		})
	}

}
