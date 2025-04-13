/*
Copyright © 2025 Jose Clavero Anderica (jose.clavero.anderica@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package utils

import (
	"fmt"
)

type DebugMode int

// DebugMode defines different levels of debug output: INFO (INFO), WARN, and ERROR.
// Each mode is associated with a specific color for terminal output.
//
// The print() method returns a colored string label representing the mode,
// which is used by the Debug function to format debug messages with consistent styling.
//
// Example output:
// [INFO]  Application started
// [WARN]  Cache miss
// [ERROR] Failed to connect to database
const (
	INFO DebugMode = iota
	ERROR
	WARN
)

var (
	Verbose bool
	reset   = "\033[0m"
	red     = "\033[31m"
	yellow  = "\033[33m"
	green   = "\033[32m"
)

func (mode DebugMode) print() string {
	var str string
	var color string
	if mode == INFO {
		color = green
		str = "INFO"
	} else if mode == WARN {
		color = yellow
		str = "WARN"
	} else if mode == ERROR {
		color = red
		str = "ERROR"
	}
	return fmt.Sprintf("[%s%s%s]", color, str, reset)
}

// Logger prints a formatted log message with the appropriate color and label
// based on the given DebugMode (LOG, WARN, ERROR). The output is only shown
// if verbose mode is enabled via SetVerbose(true).
//
// Example:
//
//	Logger(LOG, "Service started")
//	Logger(ERROR, "Failed to connect to DB")
func Logger(mode DebugMode, msg string) {
	if !Verbose && mode  != ERROR {
		fmt.Println(Verbose)
		return
	}
	fmt.Printf(" %s %s\n", mode.print(), msg)
}
