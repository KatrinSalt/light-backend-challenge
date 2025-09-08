package output

import (
	"fmt"
	"os"
)

// Println prints a message to stdout
func Println(message string) {
	fmt.Println(message)
}

// PrintlnErr prints an error message to stderr
func PrintlnErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// Printf prints a formatted message to stdout
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// PrintfErr prints a formatted error message to stderr
func PrintfErr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
