// Package testutil holds helpers shared by tests (cross-platform process stubs).
package testutil

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

// StubEcho returns a command that prints msg to stdout and exits 0 (Unix + Windows).
func StubEcho(msg string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// Avoid cmd.exe "echo" quirks (some strings parse differently); PowerShell is deterministic.
		return exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
			fmt.Sprintf("Write-Output %s", strconv.Quote(msg)))
	}
	return exec.Command("echo", msg)
}

// StubOK exits with status 0.
func StubOK() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "exit", "0")
	}
	return exec.Command("true")
}

// StubExit1 exits with status 1.
func StubExit1() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "exit", "1")
	}
	return exec.Command("false")
}
