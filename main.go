package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	prefix      = "x"
)

func run(args []string, dir string) int {
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

	result := 0
	if commandLine.NFlag() == 0 {
		result = splitByLineCount(file, dir, 1000)
	} else if *lineCount > 0 {
		result = splitByLineCount(file, dir, *lineCount)
	} else if *chunkCount > 0 {
		result = splitByChunkCount(file, dir, *chunkCount)
	} else if *byteCount > 0 {
		result = splitByByteCount(file, dir, *byteCount)
	} else {
		fmt.Fprintf(os.Stderr, "invalid flag\n")
		result = 1
	}

	return result
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

func splitByLineCount(file *os.File, dir string, lineCount int) int {
	scanner := bufio.NewScanner(file)
	filename := prefix + "aa"

	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
		return 1
	}

	l := lineCount
	for scanner.Scan() {
		f.WriteString(scanner.Text() + "\n")
		l--
		if l == 0 {
			f.Close()
			filename = incrementString(filename)
			f, err = os.Create(filepath.Join(dir, filename))
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
				return 1
			}
			l = lineCount
		}
	}

	f.Close()

	return 0
}

func splitByChunkCount(file *os.File, dir string, chunkCount int) int {
	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine file size\n")
		return 1
	}

	size := fi.Size() / int64(chunkCount)

	for i, filename := 0, prefix+"aa"; i < chunkCount; i, filename = i+1, incrementString(filename) {
		f, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
			return 1
		}

		// 最後のチャンクでは残り全て読み込む
		if i == chunkCount-1 {
			size = fi.Size() - size*int64(i)
		}

		buffer := make([]byte, size)
		file.Read(buffer)
		f.Write(buffer)
		f.Close()
	}

	return 0
}

func splitByByteCount(file *os.File, dir string, byteCount int) int {
	reader := bufio.NewReader(file)

	for filename := prefix + "aa"; ; filename = incrementString(filename) {
		buffer := make([]byte, byteCount)
		_, err := reader.Read(buffer)
		if err != nil {
			break
		}

		f, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
			return 1
		}

		f.Write(buffer)
		f.Close()
	}

	return 0
}

func main() {
	dir, _ := os.Getwd()
	os.Exit(run(os.Args[1:], dir))
}
