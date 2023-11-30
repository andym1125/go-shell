package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/andym1125/go-shell/builtins"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	initHistory()
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func HandleInput(w io.Writer, exit chan<- struct{}, args ...string) error {
	return handleInput(w, strings.Join(args, " "), exit)
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]
	var err error

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		err = builtins.ChangeDirectory(args...)
	case "env":
		err = builtins.EnvironmentVariables(w, args...)
	case "pwd":
		err = builtins.PrintWorkingDir(w, args...) //bash versions
	case "history":
		err = builtins.History(w, args...) //csh version
	case "echo":
		err = builtins.Echo(w, args...) //csh version
	case "repeat":
		err = builtins.RepeatCommand(w, exit, handleInput, args...) //csh version
	case "export":
		err = builtins.Export(w, args...) //zsh version
	case "sysrun":
		err = executeCommand(args[0], args[1:]...) //convenience builtin to run commands directly on system
	case "exit":
		exit <- struct{}{}
		err = nil
	default:
		err = executeCommand(name, args...)
	}

	werr := writeHistory(name + " " + strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("%w: %w", err, werr)
	}
	return werr
}

func initHistory() {
	histfile := os.Getenv("HISTFILE")
	if histfile == "" {
		os.Setenv("HISTFILE", "./.shell/history.txt")
	}
}

func writeHistory(input string) error {
	histfile := os.Getenv("HISTFILE")
	if histfile == "" {
		os.Setenv("HISTFILE", "./.shell/history.txt")
		histfile = "./.shell/history.txt"
	}

	// Open the file for appending.
	f, err := os.OpenFile(histfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer f.Close()

	// Write the input to the file.
	_, err = fmt.Fprintln(f, input)

	return err
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
