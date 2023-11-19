package utils

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsStdinTerminal returns true when the os.Stdin is a terminal (including cygwin/msys2 terminals)
func IsStdinTerminal() bool {
	fd := uintptr(os.Stdin.Fd())
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
