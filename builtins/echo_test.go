package builtins_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/andym1125/go-shell/builtins"
)

func TestEcho(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  error
	}{
		{
			name:     "test echo with no args",
			wantText: "\n",
		},
		{
			name: "test echo with args",
			args: args{
				args: []string{"Hello", "World"},
			},
			wantText: "Hello World\n",
		},
		{
			name: "test echo no newline",
			args: args{
				args: []string{"-n", "Hello", "World"},
			},
			wantText: "Hello World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			w := io.Writer(&bytes.Buffer{})

			// testing
			if err := builtins.Echo(w, tt.args.args...); tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Echo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else if err != nil {
				t.Fatalf("Echo() unexpected error: %v", err)
			}

			// "happy" path
			if gotText := w.(*bytes.Buffer).String(); gotText != tt.wantText {
				t.Fatalf("Failed to echo: got ;%v;, want ;%v;", gotText, tt.wantText)
			}
		})
	}
}
