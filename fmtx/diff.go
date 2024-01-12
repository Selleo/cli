package fmtx

import (
	"fmt"
	"io"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/wzshiming/ctc"
)

func ContentDiff(w io.Writer, a string, b string, inline bool) {
	if inline {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(a), string(b), false)
		fmt.Fprint(w, dmp.DiffPrettyText(diffs))
	} else {
		edits := myers.ComputeEdits(span.URIFromPath("A"), string(a), string(b))
		diff := fmt.Sprint(gotextdiff.ToUnified("A", "B", string(a), edits))

		lines := strings.Split(diff, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "+") {
				fmt.Fprint(w, ctc.ForegroundGreen, line, ctc.Reset, "\n")
			} else if strings.HasPrefix(line, "-") {
				fmt.Fprint(w, ctc.ForegroundRed, line, ctc.Reset, "\n")
			} else {
				fmt.Fprint(w, "", line, "\n")
			}
		}
	}

}
