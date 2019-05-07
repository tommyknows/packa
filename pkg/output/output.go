package output

import (
	"fmt"
	"os"

	a "github.com/logrusorgru/aurora"
)

// Info prints an info string to the terminal
func Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Success prints a string as a success message
func Success(format string, args ...interface{}) {
	fmt.Println(a.Green(a.Bold(fmt.Sprintf(format, args...))))
}

// Warn prints a string as a warning to the terminal
func Warn(format string, args ...interface{}) {
	fmt.Println(a.Bold(a.Yellow(fmt.Sprintf(format, args...))))
}

// Error prints an error
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, a.Bold(a.Red(fmt.Sprintf(format, args...))).String())
}
