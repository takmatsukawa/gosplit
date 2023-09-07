package main

import (
	"flag"
	"fmt"
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
				t.Errorf("Unexpected file xaa")
			}
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
				name:             "1 line",
				content:          "a",
				splitCount:       1,
				expectedContents: []string{"a"},
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

			for i, filename := 0, "xaa"; i < len(tc.expectedContents); i, filename = i+1, incrementString(filename) {
				outputFile, err := os.Open(filepath.Join(dir, filename))
				if err != nil {
					t.Errorf("Expected file %s, got error %v", filename, err)
					continue
				}
				content, err := io.ReadAll(outputFile)
				if err != nil {
					t.Errorf("Expected file %s to be readable, got error %v", filename, err)
					continue
				}
				if string(content) != tc.expectedContents[i] {
					t.Errorf("Expected %s in file %s, got %s", tc.expectedContents[i], filename, string(content))
				}
				outputFile.Close()
			}
		}
	})

	// tempDir := t.TempDir()

	// inputFilePath := filepath.Join(tempDir, "test.txt")
	// inputFile, _ := os.Create(inputFilePath)
	// for i := 0; i < 100; i++ {
	// 	inputFile.WriteString(fmt.Sprintf("Line %d\n", i+1))
	// }
	// inputFile.Close()

	// inputFile, _ = os.Open(inputFilePath)

	// splitByLineCount(inputFile, tempDir, 10)

	// for i, filename := 0, "xaa"; i < 10; i, filename = i+1, incrementString(filename) {
	// 	outputFile, err := os.Open(filepath.Join(tempDir, filename))
	// 	if err != nil {
	// 		t.Errorf("Expected file %s, got error %v", filename, err)
	// 		continue
	// 	}
	// 	content, err := io.ReadAll(outputFile)
	// 	if err != nil {
	// 		t.Errorf("Expected file %s to be readable, got error %v", filename, err)
	// 		continue
	// 	}
	// 	lines := strings.Split(string(content), "\n")
	// 	if len(lines) != 11 { // 10 lines + 1 empty line at the end
	// 		t.Errorf("Expected 10 lines in file %s, got %d: %s", filename, len(lines)-1, string(content))
	// 	}
	// 	outputFile.Close()
	// }
}

func TestSplitByChunkCount(t *testing.T) {
	tempDir := t.TempDir()

	inputFilePath := filepath.Join(tempDir, "test.txt")
	inputFile, _ := os.Create(inputFilePath)
	for i := 0; i < 10; i++ {
		inputFile.WriteString(fmt.Sprintf("%04d\n", i+1))
	}
	inputFile.Close()

	inputFile, _ = os.Open(inputFilePath)

	splitByChunkCount(inputFile, tempDir, 10)

	for i, filename := 0, "xaa"; i < 10; i, filename = i+1, incrementString(filename) {
		outputFile, err := os.Open(filepath.Join(tempDir, filename))
		if err != nil {
			t.Errorf("Expected file %s, got error %v", filename, err)
			continue
		}
		content, err := io.ReadAll(outputFile)
		if err != nil {
			t.Errorf("Expected file %s to be readable, got error %v", filename, err)
			continue
		}
		if len(string(content)) != 5 {
			t.Errorf("Expected 5 length in file %s, got %d: %s", filename, len(string(content)), string(content))
		}
		outputFile.Close()
	}
}

func TestSplitByByteCount(t *testing.T) {
	tempDir := t.TempDir()

	inputFilePath := filepath.Join(tempDir, "test.txt")
	inputFile, _ := os.Create(inputFilePath)
	for i := 0; i < 10; i++ {
		inputFile.WriteString(fmt.Sprintf("%09d\n", i+1))
	}
	inputFile.Close()

	inputFile, _ = os.Open(inputFilePath)

	splitByByteCount(inputFile, tempDir, 10)

	for i, filename := 0, "xaa"; i < 10; i, filename = i+1, incrementString(filename) {
		outputFile, err := os.Open(filepath.Join(tempDir, filename))
		if err != nil {
			t.Errorf("Expected file %s, got error %v", filename, err)
			continue
		}
		content, err := io.ReadAll(outputFile)
		if err != nil {
			t.Errorf("Expected file %s to be readable, got error %v", filename, err)
			continue
		}
		if len(string(content)) != 10 {
			t.Errorf("Expected 10 length in file %s, got %d: %s", filename, len(string(content)), string(content))
		}
		outputFile.Close()
	}
}
