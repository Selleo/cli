package shellcmd

import (
	"fmt"
	"io"
	"strings"
)

func KeyValueToExports(w io.Writer, kv map[string]string) {
	for k, v := range kv {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		fmt.Fprint(w, "export ", k, "=", `"`, escaped, `"`, "\n")
	}
}
