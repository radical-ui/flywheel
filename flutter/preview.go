package flutter

import (
	"os"
	"os/exec"
)

func (self *Flutter) Preview() error {
	cmd := exec.Command("flutter", "run")
	cmd.Dir = self.dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	// TODO send sigint when CLI gets SIGINT

	return cmd.Run()
}
