package builtins

import (
	"io"
	"os"
	"strings"
)

func Export(w io.Writer, args ...string) error {

	keyval := strings.Split(args[0], "=")
	if len(keyval) != 2 {
		return ErrInvalidArgCount
	}

	return os.Setenv(keyval[0], keyval[1])
}
