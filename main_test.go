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
		{name: "without args", args: []string{}, want: 0},
		{name: "lとnとbは同時に指定できない", args: []string{"-l", "1", "-n", "1", "-b", "1"}, want: 1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			if got := run(tt.args, t.TempDir()); got != tt.want {
				t.Errorf("%v: run() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
