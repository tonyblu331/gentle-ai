package upgrade

import (
	"os/exec"

	"github.com/gentleman-programming/gentle-ai/internal/testutil"
)

func testExecNoop() *exec.Cmd              { return testutil.StubExit0() }
func testExecFail() *exec.Cmd              { return testutil.StubExit1() }
func testExecEcho(msg ...string) *exec.Cmd { return testutil.StubEcho(msg...) }
