// Package fortune wraps a list of clippings.Clipping with a type to allow
// them to be written in a nice format.
package fortune

import (
	"fmt"
	"io"
	"strings"

	"hawx.me/code/kindle-clippings/clippings"
)

// This is basically a worse version of
//  https://github.com/lucaswiman/stuff/blob/master/html/kindle-highlights-fortune-file/kindle-to-fortune.py

var delim = []byte("\n%\n")

type Fortunes []clippings.Clipping

func (f Fortunes) WriteTo(w io.Writer) error {
	for i, fortune := range f {
		if i > 0 {
			_, err := w.Write(delim)
			if err != nil {
				return err
			}
		}

		if err := format(w, fortune); err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\n"))
	return err
}

func format(w io.Writer, c clippings.Clipping) error {
	_, err := fmt.Fprintf(w, `
%s

    %s
    By %s
`, justify(c.Content, 80), splitTitle(c.Title), c.Author)
	return err
}

func justify(text string, numberOfChars int) string {
	words := strings.Fields(text)
	var lines []string

	curLine := "  "
	for _, word := range words {
		if len(curLine+word) > numberOfChars {
			lines = append(lines, curLine)
			curLine = "  " + word + " "
		} else {
			curLine += word + " "
		}
	}

	if len(curLine) > 0 {
		lines = append(lines, curLine)
	}

	return strings.Join(lines, "\n")
}

func splitTitle(title string) string {
	if i := strings.Index(title, ":"); i >= 0 {
		return title[:i+1] + "\n       " + title[i+1:]
	}

	return title
}
