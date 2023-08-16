package main

import (
	"flag"
	"os"
	"testing"
)

func TestFlagVar(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		// {name: "test1", args: []string{"-name", "foo"}, want: 0},
		// {name: "test2", args: []string{"-name", "foo", "-max", "1000"}, want: 1},
		// {name: "test3", args: []string{"-name", "", "-max", "123"}, want: 1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			if got := run(tt.args); got != tt.want {
				t.Errorf("%v: run() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
