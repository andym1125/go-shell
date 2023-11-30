// Author's Note: There is a very technical difference between `pwd -L` and `pwd -P`.
// This is the author's attempt to represent that difference in Go without unnecessarily farming
// out work to an Exec call.
// Behavior:
// `pwd -L`: Possibly contains symbolic links. Returns the contents of $PWD env var.
// `pwd -P`: Does not contain symbolic links. Returns value after running /bin/pwd.
package builtins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func PrintWorkingDir(w io.Writer, args ...string) error {

	//No symbolic links
	if len(args) >= 1 && args[0] == "-P" {
		res, err := exec.Command("/bin/pwd").Output()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s", res)
		return err
		//Symlinks allowed
	} else {
		fmt.Fprintln(w, os.Getenv("PWD"))
		return nil
	}
}
