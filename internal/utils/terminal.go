package utils

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsStdinTerminal returns true when the os.Stdin is a terminal (including cygwin/msys2 terminals)
func IsStdinTerminal() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}
