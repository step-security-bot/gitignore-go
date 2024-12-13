package gitignore_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"git.sr.ht/~jamesponddotco/gitignore-go"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		wantErr   error
		wantMatch string
	}{
		{
			name: "Valid gitignore file",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = filepath.Join(dir, ".gitignore")
				)

				err := os.WriteFile(path, []byte("*.log\n"), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "Empty gitignore file",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = filepath.Join(dir, ".gitignore")
				)

				err := os.WriteFile(path, []byte(""), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr: nil,
		},
		{
			name: "Non-existent file",
			setup: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "Permission denied",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = filepath.Join(dir, ".gitignore")
				)

				err := os.WriteFile(path, []byte("*.log\n"), 0o000)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr: os.ErrPermission,
		},
		{
			name: "Relative path",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = ".gitignore"
				)

				err := os.Chdir(dir)
				if err != nil {
					t.Fatalf("failed to change directory: %v", err)
				}

				err = os.WriteFile(path, []byte("*.log\n"), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "File with special characters in path",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = filepath.Join(dir, "special!@#$%.gitignore")
				)

				err := os.WriteFile(path, []byte("*.log\n"), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "File with invalid pattern",
			setup: func(t *testing.T) string {
				t.Helper()

				var (
					dir  = t.TempDir()
					path = filepath.Join(dir, ".gitignore")
				)

				err := os.WriteFile(path, []byte("[invalid-regex\n"), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				return path
			},
			wantErr: gitignore.ErrRegexCompile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.setup(t)

			file, err := gitignore.New(path)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("New(%q) = nil error, want error containing %q", path, tt.wantErr)

					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("New(%q) error = %v, want error containing %q", path, err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("New(%q) unexpected error: %v", path, err)

				return
			}

			if tt.wantMatch != "" {
				if !file.Match(tt.wantMatch) {
					t.Errorf("New(%q) created matcher failed to match %q", path, tt.wantMatch)
				}
			}
		})
	}
}

func TestNewFromLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		giveLines []string
		wantErr   error
		wantMatch string
	}{
		{
			name: "Valid single line",
			giveLines: []string{
				"*.log",
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name:      "Empty lines slice",
			giveLines: make([]string, 0),
			wantErr:   nil,
		},
		{
			name: "Multiple valid lines",
			giveLines: []string{
				"*.log",
				"*.tmp",
				"*.cache",
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "Lines with comments",
			giveLines: []string{
				"# This is a comment",
				"*.log",
				"# Another comment",
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "Lines with empty strings",
			giveLines: []string{
				"",
				"*.log",
				"",
			},
			wantErr:   nil,
			wantMatch: "test.log",
		},
		{
			name: "Invalid regex pattern",
			giveLines: []string{
				"[invalid-regex",
			},
			wantErr: gitignore.ErrRegexCompile,
		},
		{
			name: "Mix of valid and invalid patterns",
			giveLines: []string{
				"*.log",
				"[invalid-regex",
				"*.tmp",
			},
			wantErr: gitignore.ErrRegexCompile,
		},
		{
			name:      "Nil lines slice",
			giveLines: nil,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := gitignore.NewFromLines(tt.giveLines)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewFromLines(%v) = nil error, want error containing %q", tt.giveLines, tt.wantErr)

					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewFromLines(%v) error = %v, want error containing %q", tt.giveLines, err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("NewFromLines(%v) unexpected error: %v", tt.giveLines, err)

				return
			}

			if tt.wantMatch != "" {
				if !file.Match(tt.wantMatch) {
					t.Errorf("NewFromLines(%v) created matcher failed to match %q", tt.giveLines, tt.wantMatch)
				}
			}
		})
	}
}

func TestFile_Match(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		giveRule  string
		givePath  string
		wantMatch bool
	}{
		{
			name:      "Simple Match",
			giveRule:  "foo",
			givePath:  "foo",
			wantMatch: true,
		},
		{
			name:      "Simple Non-Match",
			giveRule:  "foo",
			givePath:  "bar",
			wantMatch: false,
		},
		{
			name:      "Simple Directory Match with Trailing Slash",
			giveRule:  "foo/",
			givePath:  "foo/",
			wantMatch: true,
		},
		{
			name:      "Simple Directory Match with Trailing Slash and Subdirectory",
			giveRule:  "foo/",
			givePath:  "foo/bar",
			wantMatch: true,
		},
		{
			name:      "Simple Directory Match without Trailing Slash",
			giveRule:  "foo/",
			givePath:  "foo",
			wantMatch: false,
		},
		{
			name:      "Anywhere Match",
			giveRule:  "**/foo",
			givePath:  "foo",
			wantMatch: true,
		},
		{
			name:      "Anywhere Match in Subdirectory",
			giveRule:  "**/foo",
			givePath:  "bar/foo",
			wantMatch: true,
		},
		{
			name:      "Anywhere Non-Match",
			giveRule:  "**/foo",
			givePath:  "bar/baz",
			wantMatch: false,
		},
		{
			name:      "Anywhere From Root Match",
			giveRule:  "/**/foo",
			givePath:  "foo",
			wantMatch: true,
		},
		{
			name:      "Anywhere From Root Match with Leading Slash",
			giveRule:  "/**/foo",
			givePath:  "/foo",
			wantMatch: true,
		},
		{
			name:      "Anywhere From Root Match in Subdirectory",
			giveRule:  "/**/foo",
			givePath:  "bar/foo",
			wantMatch: true,
		},
		{
			name:      "Root Extension Only Match",
			giveRule:  "/.js",
			givePath:  ".js",
			wantMatch: true,
		},
		{
			name:      "Root Extension Only Match with Trailing Slash",
			giveRule:  "/.js",
			givePath:  ".js/",
			wantMatch: true,
		},
		{
			name:      "Root Extension Only Match with Subdirectory",
			giveRule:  "/.js",
			givePath:  ".js/a",
			wantMatch: true,
		},
		{
			name:      "Root Extension Only Non-Match",
			giveRule:  "/.js",
			givePath:  ".jsa",
			wantMatch: false,
		},
		{
			name:      "Root Extension Match",
			giveRule:  "/*.js",
			givePath:  ".js",
			wantMatch: true,
		},
		{
			name:      "Root Extension Match with Trailing Slash",
			giveRule:  "/*.js",
			givePath:  ".js/",
			wantMatch: true,
		},
		{
			name:      "Root Extension Match with Subdirectory",
			giveRule:  "/*.js",
			givePath:  ".js/a",
			wantMatch: true,
		},
		{
			name:      "Root Extension Match with Subdirectory and Extension",
			giveRule:  "/*.js",
			givePath:  "a.js/a",
			wantMatch: true,
		},
		{
			name:      "Root Extension Match with Subdirectory and Extension 2",
			giveRule:  "/*.js",
			givePath:  "a.js/a.js",
			wantMatch: true,
		},
		{
			name:      "Root Extension Non-Match",
			giveRule:  "/*.js",
			givePath:  ".jsa",
			wantMatch: false,
		},
		{
			name:      "Extension Match",
			giveRule:  "*.js",
			givePath:  ".js",
			wantMatch: true,
		},
		{
			name:      "Extension Match with Trailing Slash",
			giveRule:  "*.js",
			givePath:  ".js/",
			wantMatch: true,
		},
		{
			name:      "Extension Match with Subdirectory",
			giveRule:  "*.js",
			givePath:  ".js/a",
			wantMatch: true,
		},
		{
			name:      "Extension Match with Subdirectory and Extension",
			giveRule:  "*.js",
			givePath:  "a.js/a",
			wantMatch: true,
		},
		{
			name:      "Extension Match with Subdirectory and Extension 2",
			giveRule:  "*.js",
			givePath:  "a.js/a.js",
			wantMatch: true,
		},
		{
			name:      "Extension Match with Leading Slash",
			giveRule:  "*.js",
			givePath:  "/.js",
			wantMatch: true,
		},
		{
			name:      "Extension Non-Match",
			giveRule:  "*.js",
			givePath:  ".jsa",
			wantMatch: false,
		},
		{
			name:      "Star Extension Match",
			giveRule:  ".js*",
			givePath:  ".js",
			wantMatch: true,
		},
		{
			name:      "Star Extension Match with Trailing Slash",
			giveRule:  ".js*",
			givePath:  ".js/",
			wantMatch: true,
		},
		{
			name:      "Star Extension Match with Subdirectory",
			giveRule:  ".js*",
			givePath:  ".js/a",
			wantMatch: true,
		},
		{
			name:      "Star Extension Non-Match with Subdirectory and Extension",
			giveRule:  ".js*",
			givePath:  "a.js/a",
			wantMatch: false,
		},
		{
			name:      "Star Extension Non-Match with Subdirectory and Extension 2",
			giveRule:  ".js*",
			givePath:  "a.js/a.js",
			wantMatch: false,
		},
		{
			name:      "Star Extension Match with Leading Slash",
			giveRule:  ".js*",
			givePath:  "/.js",
			wantMatch: true,
		},
		{
			name:      "Star Extension Match",
			giveRule:  ".js*",
			givePath:  ".jsa",
			wantMatch: true,
		},
		{
			name:      "Double Star Directory Match",
			giveRule:  "foo/**/",
			givePath:  "foo/",
			wantMatch: true,
		},
		{
			name:      "Double Star Directory Match with Subdirectory",
			giveRule:  "foo/**/",
			givePath:  "foo/abc/",
			wantMatch: true,
		},
		{
			name:      "Double Star Directory Match with Deep Subdirectory",
			giveRule:  "foo/**/",
			givePath:  "foo/x/y/z/",
			wantMatch: true,
		},
		{
			name:      "Double Star Directory Non-Match",
			giveRule:  "foo/**/",
			givePath:  "foo",
			wantMatch: false,
		},
		{
			name:      "Double Star Directory Non-Match with Leading Slash",
			giveRule:  "foo/**/",
			givePath:  "/foo",
			wantMatch: false,
		},
		{
			name:      "Stars with Extension Non-Match",
			giveRule:  "foo/**/*.bar",
			givePath:  "foo/",
			wantMatch: false,
		},
		{
			name:      "Stars with Extension Non-Match 2",
			giveRule:  "foo/**/*.bar",
			givePath:  "abc.bar",
			wantMatch: false,
		},
		{
			name:      "Stars with Extension Match",
			giveRule:  "foo/**/*.bar",
			givePath:  "foo/abc.bar",
			wantMatch: true,
		},
		{
			name:      "Stars with Extension Match with Trailing Slash",
			giveRule:  "foo/**/*.bar",
			givePath:  "foo/abc.bar/",
			wantMatch: true,
		},
		{
			name:      "Stars with Extension Match with Deep Subdirectory",
			giveRule:  "foo/**/*.bar",
			givePath:  "foo/x/y/z.bar",
			wantMatch: true,
		},
		{
			name:      "Stars with Extension Match with Deep Subdirectory and Trailing Slash",
			giveRule:  "foo/**/*.bar",
			givePath:  "foo/x/y/z.bar/",
			wantMatch: true,
		},
		{
			name:      "Comment Handling",
			giveRule:  "#abc",
			givePath:  "#abc",
			wantMatch: false,
		},
		{
			name:      "Escaped Comment Handling",
			giveRule:  `\#abc`,
			givePath:  "#abc",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 2",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 3",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 4",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 5",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 6",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 7",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 8",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 9",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 10",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 11",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 12",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 13",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 14",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 15",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 16",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 17",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 18",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 19",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 20",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 21",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 22",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 23",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 24",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 25",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 26",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 27",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 28",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 29",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 30",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 31",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 32",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 33",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 34",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 35",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 36",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 37",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 38",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 39",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 40",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 41",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 42",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 43",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 44",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 45",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 46",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 47",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 48",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 49",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 50",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 51",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 52",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 53",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 54",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 55",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 56",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 57",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 58",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 59",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 60",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 61",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 62",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 63",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 64",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 65",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 66",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 67",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 68",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 69",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 70",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 71",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 72",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 73",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 74",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 75",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 76",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 77",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 78",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 79",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 80",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 81",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 82",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 83",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 84",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 85",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 86",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 87",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 88",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 89",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 90",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 91",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 92",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 93",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 94",
			giveRule:  "abc",
			givePath:  "abc/a.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 95",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
		{
			name:      "Negated Pattern Handling - Should Match 96",
			giveRule:  "abc",
			givePath:  "abc/b/b.js",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			matcher, err := gitignore.NewFromLines([]string{tt.giveRule})
			if err != nil {
				t.Fatalf("failed to create matcher: %v", err)
			}

			got := matcher.Match(tt.givePath)
			if got != tt.wantMatch {
				t.Errorf("Match(%q) with giveRule %q = %v, want %v", tt.givePath, tt.giveRule, got, tt.wantMatch)
			}
		})
	}
}

func TestFile_Match_Negation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		giveRules []string
		givePath  string
		wantMatch bool
	}{
		{
			name: "Basic Negation Override",
			giveRules: []string{
				"*.log",
				"!important.log",
			},
			givePath:  "important.log",
			wantMatch: false,
		},
		{
			name: "Multiple Patterns with Negation",
			giveRules: []string{
				"temp/*",
				"*.log",
				"!temp/special.log",
			},
			givePath:  "temp/special.log",
			wantMatch: false,
		},
		{
			name: "Complex Pattern with Negation",
			giveRules: []string{
				"**/logs/**",
				"!**/logs/keep/**",
			},
			givePath:  "server/logs/keep/debug.log",
			wantMatch: false,
		},
		{
			name: "Multiple Matches with Final Negation",
			giveRules: []string{
				"*.js",
				"lib/*.js",
				"src/*.js",
				"!src/main.js",
			},
			givePath:  "src/main.js",
			wantMatch: false,
		},
		{
			name: "Negation Between Matches",
			giveRules: []string{
				"*.txt",
				"!important/*.txt",
				"important/temp/*.txt",
			},
			givePath:  "important/temp/data.txt",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			matcher, err := gitignore.NewFromLines(tt.giveRules)
			if err != nil {
				t.Fatalf("failed to create matcher: %v", err)
			}

			got := matcher.Match(tt.givePath)
			if got != tt.wantMatch {
				t.Errorf("Match(%q) = %v, want %v", tt.givePath, got, tt.wantMatch)
			}
		})
	}
}
