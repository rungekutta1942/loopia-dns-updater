package os_util

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

// ExecuteCmd Executes a command and waits for it to finish, at most 60 seconds.
// The output is logged if the log level is debug.
func ExecuteCmd(cmd string, dry bool) {

	if dry {
		slog.Info(fmt.Sprintf("Would run '%s'", cmd))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmdSplit := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(cmdSplit) == 1 {
		c = exec.CommandContext(ctx, cmdSplit[0])
	} else {
		c = exec.CommandContext(ctx, cmdSplit[0], cmdSplit[1:]...)
	}
	out, err := c.CombinedOutput()
	if err != nil {
		slog.Error(fmt.Sprintf("Can't execute hook '%s': %s", cmd, err))
	}
	slog.Debug(fmt.Sprintf("Run '%s' with output:\n%s", cmd, string(out)))
}
