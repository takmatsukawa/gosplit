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
	lineCount := commandLine.Int("l", 0, "Create split files line_count lines in length.")
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

	// lとnとbは同時に指定できない
	if commandLine.NFlag() > 1 {
		fmt.Fprintf(os.Stderr, "cannot split in more than one way\n")
		return 1
	}

	if commandLine.NFlag() == 0 {
		splitByLineCount(file, 1000)
	} else if *lineCount > 0 {
		splitByLineCount(file, *lineCount)
	} else if *chunkCount > 0 {
		splitByChunkCount(file, *chunkCount)
	} else if *byteCount > 0 {
		splitByByteCount(file, *byteCount)
	}

	_ = chunkCount
	_ = byteCount

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

func splitByLineCount(file *os.File, lineCount int) {
	scanner := bufio.NewScanner(file)
	filename := prefix + "aa"

	f, _ := os.Create(filename)
	defer f.Close()

	l := lineCount
	for scanner.Scan() {
		f.WriteString(scanner.Text() + "\n")
		l--
		if l == 0 {
			filename = incrementString(filename)
			f, _ = os.Create(filename)
			defer f.Close()
			l = lineCount
		}
	}
}

func splitByChunkCount(file *os.File, chunkCount int) {
	fi, _ := file.Stat()

	size := fi.Size() / int64(chunkCount)

	for i, filename := 0, prefix+"aa"; i < chunkCount; i, filename = i+1, incrementString(filename) {
		f, _ := os.Create(filename)
		defer f.Close()

		// 最後のチャンクでは残り全て読み込む
		if i == chunkCount-1 {
			size = fi.Size() - size*int64(i)
		}

		buffer := make([]byte, size)
		file.Read(buffer)
		f.Write(buffer)
	}
}

func splitByByteCount(file *os.File, byteCount int) {
	reader := bufio.NewReader(file)
	filename := prefix + "aa"
	buffer := make([]byte, byteCount)

	f, _ := os.Create(filename)
	defer f.Close()

	for {
		_, err := reader.Read(buffer)
		if err != nil {
			break
		}

		f.Write(buffer)
		filename = incrementString(filename)
		f, _ = os.Create(filename)
		defer f.Close()
	}
}

func main() {
	os.Exit(run(os.Args[1:]))
}
