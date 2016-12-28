package fortune

import (
	"bytes"
	"testing"

	"hawx.me/code/kindle-clippings/clippings"
)

func TestTitleWithColon(t *testing.T) {
	items := []clippings.Clipping{{
		Title:   "Something: the something something",
		Author:  "what?",
		Content: "oh.",
	}}

	var buf bytes.Buffer
	Fortunes(items).WriteTo(&buf)

	expected := "\n" +
		"  oh. \n" +
		`
    Something:
        the something something
    By what?

`

	if buf.String() != expected {
		t.Fatalf(`Expected:

%s

But got:

%s`, expected, buf.String())
	}
}
