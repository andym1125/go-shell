package builtins

import (
	"io"
	"strconv"
	"strings"
)

func RepeatCommand(w io.Writer, exit chan<- struct{}, callback func(io.Writer, string, chan<- struct{}) error, args ...string) error {
	allargs := strings.Join(args[1:], " ")
	num, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	for i := 0; i < num; i++ {
		if err = callback(w, allargs, exit); err != nil {
			return err
		}
	}

	return nil
}
