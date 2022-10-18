package clitools

import (
	"os"
	"testing"
)

type testVars struct {
	description string
	args        []string
	expect      command
	check       func(a, b command) bool
}

func TestPasreCommand(t *testing.T) {
	var tests = []testVars{
		{
			description: "Action not in lowercase",
			args: []string{
				"hTtP", "--port", "9090",
			}, expect: command{
				Action: "http",
			},
			check: func(exp, curr command) bool {
				return exp.Action == curr.Action
			},
		},
		{
			description: "Port value is not match",
			args: []string{
				"grpc", "--port", "9091",
			}, expect: command{
				Port: 9091,
			},
			check: func(exp, curr command) bool {
				return exp.Port == curr.Port
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"<student_aggregator>"}, test.args...)
		current := PasreCommand()
		if !test.check(test.expect, current) {
			t.Log(test.description)
			t.Logf("Parsed as %v", current)
			t.Logf("Want %v", test.expect)
			t.Errorf("When cli args %v", test.args)
		}
	}

}
