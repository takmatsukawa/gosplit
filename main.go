package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileSplitter struct {
	prefix string
}

func NewFileSplitter() *FileSplitter {
	return &FileSplitter{prefix: "x"}
}

func (sp *FileSplitter) splitByLineCount(file *os.File, dir string, lineCount int) int {
	reader := bufio.NewReader(file)
	filename := sp.prefix + "aa"

	var f *os.File
	var err error

	l := lineCount
	for {
		line, readErr := reader.ReadString('\n')
		if len(line) == 0 && readErr != nil {
			if readErr != io.EOF {
				fmt.Fprintf(os.Stderr, "cannot read file: %v\n", readErr)
				return 1
			}
			break
		}
		if f == nil {
			f, err = os.Create(filepath.Join(dir, filename))
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
				return 1
			}
		}
		f.WriteString(line)
		l--
		if l == 0 {
			f.Close()
			f = nil
			filename = incrementLastChar(filename)
			l = lineCount
		}
	}

	f.Close()

	return 0
}

func (sp *FileSplitter) splitByChunkCount(file *os.File, dir string, chunkCount int) int {
	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine file size: %v\n", err)
		return 1
	}

	size := fi.Size() / int64(chunkCount)

	for i, filename := 0, sp.prefix+"aa"; i < chunkCount; i, filename = i+1, incrementLastChar(filename) {
		f, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create file: %v\n", err)
			return 1
		}

		// The last chunk reads all the rest.
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

func (sp *FileSplitter) splitByByteCount(file *os.File, dir string, byteCount int) int {
	reader := bufio.NewReader(file)

	for filename := sp.prefix + "aa"; ; filename = incrementLastChar(filename) {
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

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func run(args []string, dir string) int {
	lineCount := commandLine.Int("l", 0, "Create split files line_count lines in length.")
	chunkCount := commandLine.Int("n", 0, "Split file into chunk_count smaller files.  The first n - 1 files will be of size (size of file / chunk_count ) and the last file will contain the remaining bytes.")
	byteCount := commandLine.Int("b", 0, "Create split files byte_count bytes in length.")
	if err := commandLine.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "cannnot parse flags: %v\n", err)
		return 1
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

	sp := NewFileSplitter()

	result := 0
	switch {
	case commandLine.NFlag() > 1: // Multiple flags cannot be specified at the same time
		fmt.Fprintf(os.Stderr, "cannot split in more than one way\n")
		result = 1
	case commandLine.NFlag() == 0: // No flags specified
		result = sp.splitByLineCount(file, dir, 1000)
	case *lineCount > 0:
		result = sp.splitByLineCount(file, dir, *lineCount)
	case *chunkCount > 0:
		result = sp.splitByChunkCount(file, dir, *chunkCount)
	case *byteCount > 0:
		result = sp.splitByByteCount(file, dir, *byteCount)
	default:
		fmt.Fprintf(os.Stderr, "invalid flag\n")
		result = 1
	}

	return result
}

func incrementLastChar(s string) string {
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
	dir, _ := os.Getwd()
	os.Exit(run(os.Args[1:], dir))
}
