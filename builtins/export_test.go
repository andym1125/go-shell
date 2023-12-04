package builtins_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/andym1125/go-shell/builtins"
)

func TestExport(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name       string
		args       args
		wantKeyVal []string
		wantErr    error
	}{
		{
			name:    "test export with no args",
			wantErr: builtins.ErrInvalidArgCount,
		},
		{
			name: "test export with incorrect arg",
			args: args{
				args: []string{"HelloWorld"},
			},
			wantErr: builtins.ErrInvalidArgCount,
		},
		{
			name: "test export happy path",
			args: args{
				args: []string{"Hello=World"},
			},
			wantKeyVal: []string{"Hello", "World"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			w := io.Writer(&bytes.Buffer{})

			// testing
			if err := builtins.Export(w, tt.args.args...); tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Echo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else if err != nil {
				t.Fatalf("Echo() unexpected error: %v", err)
			}

			// "happy" path
			if os.Getenv(tt.wantKeyVal[0]) != tt.wantKeyVal[1] {
				t.Fatalf("Failed to echo: got ;%v;, want ;%v;", os.Getenv(tt.wantKeyVal[0]), tt.wantKeyVal[1])
			}
		})
	}
}
