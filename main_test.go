package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_handleInput(t *testing.T) {
	// t.Parallel()
	exit := make(chan struct{}, 2)
	w := &bytes.Buffer{}
	test := []string{"cd .", "ls", "env", "pwd", "history", "repeat 2 echo repeater", "sysrun echo sysrun", "echo hello", "export foo=bar"}
	for _, tt := range test {
		t.Run(tt, func(t *testing.T) {
			// t.Parallel() //Removed bc it was preventing code coverage from working for some reason
			err := handleInput(w, tt, exit)
			require.NoError(t, err)

		})
	}

	err := handleInput(w, "exit", exit)
	require.NoError(t, err)
}

func Test_runLoop(t *testing.T) {
	t.Parallel()
	os.Setenv("HISTFILE", "")
	exitCmd := strings.NewReader("exit\n")
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantW    string
		wantErrW string
	}{
		{
			name: "no error",
			args: args{
				r: exitCmd,
			},
		},
		{
			name: "read error should have no effect",
			args: args{
				r: iotest.ErrReader(io.EOF),
			},
			wantErrW: "EOF",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}
			errW := &bytes.Buffer{}

			exit := make(chan struct{}, 2)
			// run the loop for 10ms
			go runLoop(tt.args.r, w, errW, exit)
			time.Sleep(10 * time.Millisecond)
			exit <- struct{}{}

			require.NotEmpty(t, w.String())
			if tt.wantErrW != "" {
				require.Contains(t, errW.String(), tt.wantErrW)
			} else {
				require.Empty(t, errW.String())
			}
		})
	}
}

func TestMetricNoLongerUseful(t *testing.T) {
	exit := make(chan struct{}, 2)
	w := &bytes.Buffer{}
	r := strings.NewReader("exit\n")
	os.Setenv("HISTFILE", " ")

	runLoop(r, w, w, exit)
	exit <- struct{}{}

	require.Error(t, handleInput(w, "printf", make(chan struct{}, 2)))
}
