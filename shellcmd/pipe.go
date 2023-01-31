package shellcmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/wzshiming/ctc"
)

func Pipe(ctx context.Context, w io.Writer, secrets map[string]string, args []string) error {
	fmt.Fprintf(
		w,
		"%sStarting service: %s%s\n",
		ctc.ForegroundYellow,
		strings.Join(args, " "),
		ctc.Reset,
	)
	cmdName := args[0]
	cmdArgs := args[1:]
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	envs := []string{}
	envs = append(envs, os.Environ()...)
	for k, v := range secrets {
		// TODO: needs escaping
		envs = append(envs, fmt.Sprint(k, "=", v))
	}
	cmd.Env = envs
	cmd.Stdin = os.Stdin

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errOut, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer out.Close()

		reader := bufio.NewReader(out)

	LOOP:
		for {
			select {

			case <-cancelCtx.Done():
				fmt.Println("stopping stdout pipe")
				break LOOP

			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					cancel()
					return
				}
				fmt.Fprintf(w, line)
			}
		}
	}()
	go func() {
		defer errOut.Close()

		reader := bufio.NewReader(errOut)

	LOOP:
		for {
			select {

			case <-cancelCtx.Done():
				fmt.Fprintf(w, "[stderr] stopping stderr pipe\n")
				break LOOP

			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(w, "[stderr] err line reading\n")
					cancel()
					return
				}
				fmt.Fprintf(w, "%s%s%s", ctc.ForegroundRed, line, ctc.Reset)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	// TODO: add signal notify

	fmt.Fprintf(w, "Stopping\n")

	return err
}
