package pattern_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/gitignore-go/internal/pattern"
)

type errorReader struct {
	readCount int
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	r.readCount++

	if r.readCount > 1 {
		return 0, errors.New("forced read error")
	}

	return copy(p, "*.log\n"), nil
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr error
	}{
		{
			name:    "Empty input",
			input:   "",
			wantLen: 0,
		},
		{
			name:    "Single pattern",
			input:   "*.log",
			wantLen: 1,
		},
		{
			name:    "Multiple patterns",
			input:   "*.log\n*.tmp\n*.cache",
			wantLen: 3,
		},
		{
			name:    "Pattern with comments",
			input:   "# This is a comment\n*.log\n# Another comment",
			wantLen: 1,
		},
		{
			name:    "Pattern with empty lines",
			input:   "\n*.log\n\n",
			wantLen: 1,
		},
		{
			name:    "Pattern with negation",
			input:   "!*.log",
			wantLen: 1,
		},
		{
			name:    "Pattern with escaped comment",
			input:   `\#not-a-comment`,
			wantLen: 1,
		},
		{
			name:    "Pattern with escaped negation",
			input:   `\!not-negated`,
			wantLen: 1,
		},
		{
			name:    "Pattern with trailing slash",
			input:   "dir/",
			wantLen: 1,
		},
		{
			name:    "Pattern with double asterisk",
			input:   "**/foo",
			wantLen: 1,
		},
		{
			name:    "Pattern with escaped asterisk",
			input:   `\*literal-asterisk`,
			wantLen: 1,
		},
		{
			name:    "Pattern with question mark",
			input:   "file?.txt",
			wantLen: 1,
		},
		{
			name:    "Pattern with leading slash",
			input:   "/root-only.txt",
			wantLen: 1,
		},
		{
			name:    "Pattern starting with escaped hash",
			input:   `\#comment`,
			wantLen: 1,
		},
		{
			name:    "Pattern starting with escaped exclamation",
			input:   `\!important`,
			wantLen: 1,
		},
		{
			name:    "Pattern with wildcard extension without leading slash",
			input:   "foo/*.txt",
			wantLen: 1,
		},
		{
			name:    "Pattern with double asterisk at root",
			input:   "/**/foo",
			wantLen: 1,
		},
		{
			name:    "Pattern with hash after negation",
			input:   "!#literal-hash",
			wantLen: 1,
		},
		{
			name:    "Pattern with exclamation after negation",
			input:   "!!literal-exclamation",
			wantLen: 1,
		},
		{
			name:    "Invalid regex pattern",
			input:   "[invalid-regex",
			wantErr: pattern.ErrInvalidRegex,
		},
		{
			name:    "Scanner error",
			input:   "test",
			wantErr: pattern.ErrScanningFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r io.Reader
			if tt.name == "Scanner error" {
				r = &errorReader{}
			} else {
				r = strings.NewReader(tt.input)
			}

			patterns, err := pattern.Parse(r)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("Parse(%q) = nil error, want error", tt.input)
				}

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}

			if len(patterns) != tt.wantLen {
				t.Fatalf("Parse(%q) returned %d patterns, want %d", tt.input, len(patterns), tt.wantLen)
			}

			for i, p := range patterns {
				if p.Regex == nil {
					t.Errorf("Pattern[%d].Regex is nil", i)
				}
			}
		})
	}
}
