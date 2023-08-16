package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	prefix      = "x"
)

func run(args []string) int {
	lineCount := commandLine.Int("l", 1000, "Create split files line_count lines in length.")
	chunkCount := commandLine.Int("n", 0, "Split file into chunk_count smaller files.  The first n - 1 files will be of size (size of file / chunk_count ) and the last file will contain the remaining bytes.")
	byteCount := commandLine.Int("b", 0, "Create split files byte_count bytes in length.")
	if err := commandLine.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "cannnot parse flags: %v\n", err)
	}

	var file *os.File
	var err error
	if commandLine.NArg() == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(commandLine.Arg(0))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open file: %v\n", err)
		return 1
	}

	if *lineCount <= 0 {
		fmt.Fprintf(os.Stderr, "line count must be positive\n")
		return 1
	}

	_ = lineCount
	_ = chunkCount
	_ = byteCount

	scanner := bufio.NewScanner(file)
	filename := prefix + "aa"

	f, _ := os.Create(filename)
	defer f.Close()

	l := *lineCount
	for scanner.Scan() {
		f.WriteString(scanner.Text() + "\n")
		l--
		if l == 0 {
			filename = incrementString(filename)
			f, _ = os.Create(filename)
			defer f.Close()
			l = *lineCount
		}
	}

	return 0
}

func incrementString(s string) string {
	runes := []rune(s)
	carry := true

	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'z' {
			runes[i]++
			carry = false
			break
		} else {
			runes[i] = 'a'
		}
	}

	if carry {
		runes = append([]rune{'a'}, runes...)
	}

	return string(runes)
}

func main() {
	os.Exit(run(os.Args[1:]))
}
