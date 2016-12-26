package fortune

import (
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
		_, err := w.Write([]byte("\n" + justify(fortune.Content, 70)))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte("\n\n    " + fortune.Title))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("\n    By " + fortune.Author))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\n"))
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
