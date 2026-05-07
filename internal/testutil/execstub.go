// Package testutil provides cross-platform subprocess stubs for tests.
// Unix shells expose echo/true/false as binaries; Windows does not.
package testutil

import (
	"os/exec"
	"runtime"
)

// StubEcho returns a command that prints args (space-separated) and exits 0.
func StubEcho(parts ...string) *exec.Cmd {
	if len(parts) == 0 {
		return StubExit0()
	}
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", append([]string{"/c", "echo"}, parts...)...)
	}
	return exec.Command("echo", parts...)
}

// StubExit0 runs a no-op command that exits with status 0.
func StubExit0() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "exit", "0")
	}
	return exec.Command("true")
}

// StubExit1 runs a command that exits with status 1.
func StubExit1() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "exit", "1")
	}
	return exec.Command("false")
}
