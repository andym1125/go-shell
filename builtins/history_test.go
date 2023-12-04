package builtins_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/andym1125/go-shell/builtins"
)

// This is a hack to get around odd env issues with testing. Should be ignored during GH Action testing
// Feel free to change it during local testing
var ROOTED_HISTFILE = "/Users/andy/Documents/workshop/CSCE4600/Project2/builtins/historytest.txt"

func TestHistory(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  error
		unsetEnv bool
	}{
		{
			name:     "test history with no args",
			wantText: "0 0\n1 1\n2 2\n",
		},
		{
			name: "test history w both args",
			args: args{
				args: []string{"-rh", "3"},
			},
			wantText: "2\n1\n0\n",
		},
		{
			name: "test history w 2 length",
			args: args{
				args: []string{"2"},
			},
			wantText: "1 1\n2 2\n",
		},
		{
			name: "test history w reverse",
			args: args{
				args: []string{"-r", "3"},
			},
			wantText: "2 2\n1 1\n0 0\n",
		},
		{
			name: "test history w no nums",
			args: args{
				args: []string{"-h", "3"},
			},
			wantText: "0\n1\n2\n",
		},
		{
			name:     "test history with not env var",
			unsetEnv: true,
			wantErr:  errors.New("open : no such file or directory: HISTFILE not set. You can set it via `export HISTFILE=/path/to/file`"),
		},
	}

	if os.Getenv("HISTFILE") != "" {
		ROOTED_HISTFILE = os.Getenv("HISTFILE")
	}

	os.Remove(ROOTED_HISTFILE)
	if err := os.WriteFile(ROOTED_HISTFILE, []byte("0\n1\n2\n"), 0644); err != nil {
		t.Fatalf("Failed to write to history file: %v", err)
	}

	for _, tt := range tests {
		fmt.Print(os.Getwd())
		if !tt.unsetEnv {
			os.Setenv("HISTFILE", ROOTED_HISTFILE)
		} else {
			os.Unsetenv("HISTFILE")
		}
		t.Run(tt.name, func(t *testing.T) {
			// setup
			w := io.Writer(&bytes.Buffer{})

			// testing
			if err := builtins.History(w, tt.args.args...); tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("History() error = (%v), wantErr (%v)", err, tt.wantErr)
				}
				return
			} else if err != nil {
				t.Fatalf("History() unexpected error: %v", err)
			}

			// "happy" path
			if gotText := w.(*bytes.Buffer).String(); gotText != tt.wantText {
				t.Fatalf("Failed to history: got ;%v;, want ;%v;", gotText, tt.wantText)
			}
		})
	}
}
