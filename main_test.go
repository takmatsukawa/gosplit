package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFlagVar(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{name: "引数未指定", args: []string{}, want: 0},
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

func TestSplitByLineCount(t *testing.T) {

	t.Run("No file created", func(t *testing.T) {
		tests := []struct {
			name    string
			content string
		}{
			{
				name:    "Empty file",
				content: "",
			},
		}

		for _, tc := range tests {
			dir := t.TempDir()
			inputFilePath := filepath.Join(dir, "test.txt")
			inputFile, _ := os.Create(inputFilePath)
			inputFile.WriteString(tc.content)
			inputFile.Close()

			inputFile, _ = os.Open(inputFilePath)

			splitByLineCount(inputFile, dir, 10)

			if _, err := os.Stat(filepath.Join(dir, "xaa")); err == nil { // xaaが存在する
				t.Errorf("%s: Unexpected file xaa", tc.name)
			}

			inputFile.Close()
		}
	})

	t.Run("More than 1 file created", func(t *testing.T) {
		tests := []struct {
			name             string
			content          string
			splitCount       int
			expectedContents []string
		}{
			{
				name:             "A file containing 1 line into 1 file",
				content:          "a",
				splitCount:       1,
				expectedContents: []string{"a"},
			},
			{
				name:             "A file containing 2 lines into 1 file",
				content:          "a\nb",
				splitCount:       2,
				expectedContents: []string{"a\nb"},
			},
			{
				name:             "A file containing 2 lines into 2 files",
				content:          "a\nb",
				splitCount:       1,
				expectedContents: []string{"a\n", "b"},
			},
		}

		for _, tc := range tests {
			dir := t.TempDir()
			inputFilePath := filepath.Join(dir, "test.txt")
			inputFile, _ := os.Create(inputFilePath)
			inputFile.WriteString(tc.content)
			inputFile.Close()

			inputFile, _ = os.Open(inputFilePath)

			splitByLineCount(inputFile, dir, tc.splitCount)

			filename := "xaa"
			for i := 0; i < len(tc.expectedContents); i, filename = i+1, incrementString(filename) {
				outputFile, err := os.Open(filepath.Join(dir, filename))
				if err != nil {
					t.Errorf("%s: Expected file %s, got error %v", tc.name, filename, err)
					continue
				}
				content, err := io.ReadAll(outputFile)
				if err != nil {
					t.Errorf("%s: Expected file %s to be readable, got error %v", tc.name, filename, err)
					continue
				}
				if string(content) != tc.expectedContents[i] {
					t.Errorf("%s: Expected %s in file %s, got %s", tc.name, tc.expectedContents[i], filename, string(content))
				}
				outputFile.Close()
			}

			if _, err := os.Stat(filepath.Join(dir, incrementString(filename))); err == nil { // 余分なファイルが存在する
				t.Errorf("%s: Unexpected file %s", tc.name, filename)
			}

			inputFile.Close()
		}
	})
}

func TestSplitByChunkCount(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		chunkCount       int
		expectedContents []string
	}{
		{
			name:             "Empty file into 1 file",
			content:          "",
			chunkCount:       1,
			expectedContents: []string{""},
		},
		{
			name:             "Empty file into 2 file",
			content:          "",
			chunkCount:       2,
			expectedContents: []string{"", ""},
		},
		{
			name:             "A file into 1 file",
			content:          "a",
			chunkCount:       1,
			expectedContents: []string{"a"},
		},
		{
			name:             "A file into 2 file",
			content:          "ab",
			chunkCount:       2,
			expectedContents: []string{"a", "b"},
		},
	}

	for _, tc := range tests {
		dir := t.TempDir()
		inputFilePath := filepath.Join(dir, "test.txt")
		inputFile, _ := os.Create(inputFilePath)
		inputFile.WriteString(tc.content)
		inputFile.Close()

		inputFile, _ = os.Open(inputFilePath)

		splitByChunkCount(inputFile, dir, tc.chunkCount)

		filename := "xaa"
		for i := 0; i < len(tc.expectedContents); i, filename = i+1, incrementString(filename) {
			outputFile, err := os.Open(filepath.Join(dir, filename))
			if err != nil {
				t.Errorf("%s: Expected file %s, got error %v", tc.name, filename, err)
				continue
			}
			content, err := io.ReadAll(outputFile)
			if err != nil {
				t.Errorf("%s: Expected file %s to be readable, got error %v", tc.name, filename, err)
				continue
			}
			if string(content) != tc.expectedContents[i] {
				t.Errorf("%s: Expected %s in file %s, got %s", tc.name, tc.expectedContents[i], filename, string(content))
			}
			outputFile.Close()
		}

		if _, err := os.Stat(filepath.Join(dir, incrementString(filename))); err == nil { // 余分なファイルが存在する
			t.Errorf("%s: Unexpected file %s", tc.name, filename)
		}

		inputFile.Close()
	}
}

func TestSplitByByteCount(t *testing.T) {
	t.Run("No file created", func(t *testing.T) {
		tests := []struct {
			name    string
			content string
		}{
			{
				name:    "Empty file",
				content: "",
			},
		}

		for _, tc := range tests {
			dir := t.TempDir()
			inputFilePath := filepath.Join(dir, "test.txt")
			inputFile, _ := os.Create(inputFilePath)
			inputFile.WriteString(tc.content)
			inputFile.Close()

			inputFile, _ = os.Open(inputFilePath)

			splitByByteCount(inputFile, dir, 10)

			if _, err := os.Stat(filepath.Join(dir, "xaa")); err == nil { // xaaが存在する
				t.Errorf("%s: Unexpected file xaa", tc.name)
			}

			inputFile.Close()
		}
	})

	t.Run("More than 1 file created", func(t *testing.T) {
		tests := []struct {
			name             string
			content          string
			byteCount        int
			expectedContents []string
		}{
			{
				name:             "A file into 1 file",
				content:          "a",
				byteCount:        1,
				expectedContents: []string{"a"},
			},
			{
				name:             "A file into 2 file",
				content:          "ab",
				byteCount:        1,
				expectedContents: []string{"a", "b"},
			},
		}

		for _, tc := range tests {
			dir := t.TempDir()
			inputFilePath := filepath.Join(dir, "test.txt")
			inputFile, _ := os.Create(inputFilePath)
			inputFile.WriteString(tc.content)
			inputFile.Close()

			inputFile, _ = os.Open(inputFilePath)

			splitByByteCount(inputFile, dir, tc.byteCount)

			filename := "xaa"
			for i := 0; i < len(tc.expectedContents); i, filename = i+1, incrementString(filename) {
				outputFile, err := os.Open(filepath.Join(dir, filename))
				if err != nil {
					t.Errorf("%s: Expected file %s, got error %v", tc.name, filename, err)
					continue
				}
				content, err := io.ReadAll(outputFile)
				if err != nil {
					t.Errorf("%s: Expected file %s to be readable, got error %v", tc.name, filename, err)
					continue
				}
				if string(content) != tc.expectedContents[i] {
					t.Errorf("%s: Expected %s in file %s, got %s", tc.name, tc.expectedContents[i], filename, string(content))
				}
				outputFile.Close()
			}

			if _, err := os.Stat(filepath.Join(dir, incrementString(filename))); err == nil { // 余分なファイルが存在する
				t.Errorf("%s: Unexpected file %s", tc.name, filename)
			}

			inputFile.Close()
		}
	})
}

func TestIncrementString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "xaa", want: "xab"},
		{in: "xaz", want: "xba"},
		{in: "xzy", want: "xzz"},
	}
	for _, tt := range tests {
		if got := incrementString(tt.in); got != tt.want {
			t.Errorf("incrementString(\"%s\") = %v, want %v", tt.in, got, tt.want)
		}
	}
}
