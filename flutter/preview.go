package flutter

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func (self *Flutter) Preview(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "flutter", "run")
	cmd.Dir = self.dir

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Cancel = func() error {
		_, err := fmt.Fprint(stdin, "q")

		go func() {
			time.Sleep(3000)
			cmd.Process.Signal(syscall.SIGTERM)
		}()

		return err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // TODO perhaps this is unix only, consider something that would work for windows
	}

	cmd.Start()

	go func() {
		defer stdin.Close() // Ensure we close the pipe when done
		_, _ = io.Copy(stdin, os.Stdin)
	}()

	return cmd.Wait()
}
