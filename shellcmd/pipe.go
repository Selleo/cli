package shellcmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/wzshiming/ctc"
)

func Pipe(ctx context.Context, w io.Writer, secrets map[string]string, args []string) error {
	cmdName := args[0]
	cmdArgs := args[1:]
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	envs := []string{}
	envs = append(envs, os.Environ()...)
	for k, v := range secrets {
		envs = append(envs, fmt.Sprint(k, "=", v))
		fmt.Fprintf(w, "exporting %s%s%s\n", ctc.ForegroundGreen, k, ctc.Reset)
	}
	cmd.Env = envs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Fprintf(
		w,
		"%sStarting service: %s%s\n",
		ctc.ForegroundYellow,
		strings.Join(args, " "),
		ctc.Reset,
	)
	err := cmd.Run()

	return err
}
