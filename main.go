package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func run(args []string) int {
	lineCount := commandLine.Int("l", 1000, "line count")
	chunkCount := commandLine.Int("n", 0, "chunk count")
	byteCount := commandLine.Int("b", 0, "byte count")
	if err := commandLine.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "cannnot parse flags: %v\n", err)
	}

	_ = lineCount
	_ = chunkCount
	_ = byteCount

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
