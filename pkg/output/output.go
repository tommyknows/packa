package output

import (
	"bufio"
	"fmt"
	"io"
	"os"

	a "github.com/logrusorgru/aurora"
)

var stdout = io.Writer(os.Stdout)
var stderr = io.Writer(os.Stderr)

func Set(out, err io.Writer) {
	stdout = out
	stderr = err
}

// Info prints an info string to the terminal
func Info(format string, args ...interface{}) {
	fmt.Fprintf(stdout, format+"\n", args...)
}

// Success prints a string as a success message
// aka bold green
func Success(format string, args ...interface{}) {
	s := fmt.Sprintf(format+"\n", args...)
	fmt.Fprint(stdout, a.Green(s).Bold().String())
}

// Warn prints a string as a warning to the terminal
// aka bold yellow
func Warn(format string, args ...interface{}) {
	s := fmt.Sprintf(format+"\n", args...)
	fmt.Fprint(stdout, a.Yellow(s).Bold().String())
}

// Error prints an error
// aka bold red
func Error(format string, args ...interface{}) {
	s := fmt.Sprintf(format+"\n", args...)
	fmt.Fprint(stderr, a.Red(s).Bold().String())
}

// WithConfirmation prints the supplied message as an info
// and waits for confirmation of the user. The default choice
// for the confirmation, and thus for the returned boolean, is false
func WithConfirmation(format string, args ...interface{}) bool {
	Info(format, args...)
	reader := bufio.NewReader(os.Stdin)
	Info("confirm (y/N):")
	text, _ := reader.ReadString('\n')
	return text == "y\n" || text == "Y\n"
}
