package builtins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	HIST_LEN = 16
)

func History(w io.Writer, args ...string) error {

	var (
		reverse bool
		no_num  bool
		histLen int = HIST_LEN
	)

	// Parse args
	for _, arg := range args {
		switch arg {
		case "-h":
			no_num = true
		case "-r":
			reverse = true
		case "-hr", "-rh":
			reverse = true
			no_num = true
		default:
			newLen, err := strconv.Atoi(arg)
			if err == nil {
				histLen = newLen
			}
		}
	}

	file, err := os.Open(os.Getenv("HISTFILE"))
	if err != nil {
		if os.Getenv("HISTFILE") == "" {
			err = fmt.Errorf("%w: HISTFILE not set. You can set it via `export HISTFILE=/path/to/file`", err)
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	recent := []string{}
	for scanner.Scan() {
		lineCount++
		recent = append(recent, scanner.Text())
		if len(recent) > histLen {
			recent = recent[1:]
		}
	}

	if reverse {
		for i, j := 0, len(recent)-1; i < j; i, j = i+1, j-1 {
			recent[i], recent[j] = recent[j], recent[i]
		}
	}

	//Print history
	for i, line := range recent {
		count := lineCount - len(recent) + i
		if reverse {
			count = lineCount - i - 1
		}

		if !no_num {
			fmt.Fprintf(w, "%d %s\n", count, line)
		} else {
			fmt.Fprintf(w, "%s\n", line)
		}
	}

	return nil
}
