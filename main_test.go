package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFlagVar(t *testing.T) {
	tests := map[string]struct {
		args []string
		want int
	}{
		"argument not specified":                              {args: []string{}, want: 0},
		"Multiple flags cannot be specified at the same time": {args: []string{"-l", "1", "-n", "1", "-b", "1"}, want: 1},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			if got := run(tt.args, t.TempDir()); got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitByLineCount(t *testing.T) {

	sp := NewFileSplitter()

	t.Run("No file created", func(t *testing.T) {
		tests := map[string]struct {
			content string
		}{
			"Empty file": {
				content: "",
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				dir := t.TempDir()
				inputFilePath := filepath.Join(dir, "test.txt")
				inputFile, _ := os.Create(inputFilePath)
				inputFile.WriteString(tc.content)
				inputFile.Close()

				inputFile, _ = os.Open(inputFilePath)

				sp.splitByLineCount(inputFile, dir, 10)

				if _, err := os.Stat(filepath.Join(dir, "xaa")); err == nil { // xaaが存在する
					t.Errorf("Unexpected file xaa")
				}

				inputFile.Close()
			})
		}
	})

	t.Run("More than 1 file created", func(t *testing.T) {
		tests := map[string]struct {
			content          string
			splitCount       int
			expectedContents []string
		}{
			"A file containing 1 line into 1 file": {
				content:          "a",
				splitCount:       1,
				expectedContents: []string{"a"},
			},
			"A file containing 2 lines into 1 file": {
				content:          "a\nb",
				splitCount:       2,
				expectedContents: []string{"a\nb"},
			},
			"A file containing 2 lines into 2 files": {
				content:          "a\nb",
				splitCount:       1,
				expectedContents: []string{"a\n", "b"},
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				dir := t.TempDir()
				inputFilePath := filepath.Join(dir, "test.txt")
				inputFile, _ := os.Create(inputFilePath)
				inputFile.WriteString(tc.content)
				inputFile.Close()

				inputFile, _ = os.Open(inputFilePath)

				sp.splitByLineCount(inputFile, dir, tc.splitCount)

				filename := "xaa"
				for i := 0; i < len(tc.expectedContents); i, filename = i+1, incrementLastChar(filename) {
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

				if _, err := os.Stat(filepath.Join(dir, incrementLastChar(filename))); err == nil { // 余分なファイルが存在する
					t.Errorf("Unexpected file %s", filename)
				}

				inputFile.Close()
			})
		}
	})
}

func TestSplitByChunkCount(t *testing.T) {

	sp := NewFileSplitter()

	tests := map[string]struct {
		content          string
		chunkCount       int
		expectedContents []string
	}{
		"Empty file into 1 file": {
			content:          "",
			chunkCount:       1,
			expectedContents: []string{""},
		},
		"Empty file into 2 file": {
			content:          "",
			chunkCount:       2,
			expectedContents: []string{"", ""},
		},
		"A file into 1 file": {
			content:          "a",
			chunkCount:       1,
			expectedContents: []string{"a"},
		},
		"A file into 2 file": {
			content:          "ab",
			chunkCount:       2,
			expectedContents: []string{"a", "b"},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			inputFilePath := filepath.Join(dir, "test.txt")
			inputFile, _ := os.Create(inputFilePath)
			inputFile.WriteString(tc.content)
			inputFile.Close()

			inputFile, _ = os.Open(inputFilePath)

			sp.splitByChunkCount(inputFile, dir, tc.chunkCount)

			filename := "xaa"
			for i := 0; i < len(tc.expectedContents); i, filename = i+1, incrementLastChar(filename) {
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

			if _, err := os.Stat(filepath.Join(dir, incrementLastChar(filename))); err == nil { // 余分なファイルが存在する
				t.Errorf("Unexpected file %s", filename)
			}

			inputFile.Close()
		})
	}
}

func TestSplitByByteCount(t *testing.T) {

	sp := NewFileSplitter()

	t.Run("No file created", func(t *testing.T) {
		tests := map[string]struct {
			content string
		}{
			"Empty file": {
				content: "",
			},
		}

		for name, tt := range tests {
			tt := tt
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				dir := t.TempDir()
				inputFilePath := filepath.Join(dir, "test.txt")
				inputFile, _ := os.Create(inputFilePath)
				inputFile.WriteString(tt.content)
				inputFile.Close()

				inputFile, _ = os.Open(inputFilePath)

				sp.splitByByteCount(inputFile, dir, 10)

				if _, err := os.Stat(filepath.Join(dir, "xaa")); err == nil { // xaaが存在する
					t.Errorf("Unexpected file xaa")
				}

				inputFile.Close()
			})
		}
	})

	t.Run("More than 1 file created", func(t *testing.T) {
		tests := map[string]struct {
			content          string
			byteCount        int
			expectedContents []string
		}{
			"A file into 1 file": {
				content:          "a",
				byteCount:        1,
				expectedContents: []string{"a"},
			},
			"A file into 2 file": {
				content:          "ab",
				byteCount:        1,
				expectedContents: []string{"a", "b"},
			},
		}

		for name, tt := range tests {
			tt := tt
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				dir := t.TempDir()
				inputFilePath := filepath.Join(dir, "test.txt")
				inputFile, _ := os.Create(inputFilePath)
				inputFile.WriteString(tt.content)
				inputFile.Close()

				inputFile, _ = os.Open(inputFilePath)

				sp.splitByByteCount(inputFile, dir, tt.byteCount)

				filename := "xaa"
				for i := 0; i < len(tt.expectedContents); i, filename = i+1, incrementLastChar(filename) {
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
					if string(content) != tt.expectedContents[i] {
						t.Errorf("Expected %s in file %s, got %s", tt.expectedContents[i], filename, string(content))
					}
					outputFile.Close()
				}

				if _, err := os.Stat(filepath.Join(dir, incrementLastChar(filename))); err == nil { // 余分なファイルが存在する
					t.Errorf("Unexpected file %s", filename)
				}

				inputFile.Close()
			})
		}
	})
}

func TestIncrementLastChar(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "xaa", want: "xab"},
		{in: "xaz", want: "xba"},
		{in: "xzy", want: "xzz"},
	}
	for _, tt := range tests {
		if got := incrementLastChar(tt.in); got != tt.want {
			t.Errorf("incrementString(\"%s\") = %v, want %v", tt.in, got, tt.want)
		}
	}
}
