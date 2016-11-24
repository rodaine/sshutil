package sshutil

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// StdPrompter is the default prompter, writing the question to os.Stdout and
// reading user input from os.Stdin
var StdPrompter Prompter = IOPrompt(os.Stdin, os.Stdout)

// A Prompter requests a response to a specified question, usually from the
// end user.
type Prompter interface {
	// Prompt shares the question (typically from the end user in a terminal) and
	// returns the user's response. If echo is false, the written response from
	// the user should be suppressed in the input mechanism (eg, masked).
	Prompt(question string, echo bool) (response string, err error)
}

// An IOPrompter prompts the user via file descriptors. Typically, streams
// are used (such as os.Stdout and os.Stdin). IOPrompter satisfies the Prompter
// interface.
type IOPrompter struct {
	in, out *os.File
}

// IOPrompt creates an IOPrompter from the provided files
func IOPrompt(in, out *os.File) IOPrompter {
	return IOPrompter{
		in:  in,
		out: out,
	}
}

// Prompt writes the question to the outbound file, and reads the response from
// the inbound file. This function blocks until EOF is returned by the inbound
// file or a new line is reached (ie, user hits return). If echo is false, and
// the file descriptor for the inbound file is attached to a TTY, the input
// will be suppressed.
func (p IOPrompter) Prompt(question string, echo bool) (string, error) {
	fmt.Fprint(p.out, question)
	if echo {
		return p.promptEcho()
	}
	return p.promptHide()
}

func (p IOPrompter) promptEcho() (string, error) {
	s := bufio.NewScanner(p.in)
	s.Split(bufio.ScanLines)
	s.Scan()
	return s.Text(), s.Err()
}

func (p IOPrompter) promptHide() (string, error) {
	b, err := terminal.ReadPassword(int(p.in.Fd()))
	return string(b), err
}

var _ Prompter = IOPrompter{}
