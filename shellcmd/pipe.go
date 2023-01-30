package shellcmd

import (
	"context"
	"fmt"
	"io"
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
	return nil
}
