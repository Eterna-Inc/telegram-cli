package utils

import (
	"os/exec"
	"runtime"
)

func RunCmd(cmdline string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdline)
	} else {
		cmd = exec.Command("bash", "-lc", cmdline)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}
