package builtins

import (
	"fmt"
	"io"
	"strings"
)

func Echo(w io.Writer, args ...string) error {

	if len(args) == 0 {
		_, err := fmt.Fprintln(w, "")
		return err
	}

	r := args[0] == "-n"
	if r {
		args = args[1:]
	}

	_, err := fmt.Fprint(w, strings.Join(args, " "))

	if !r {
		_, err = fmt.Fprintln(w, "")
	}

	return err
}
