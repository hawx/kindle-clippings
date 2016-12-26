// Package clippings provides a Reader for the Kindle's clipping format.
package clippings

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	delim         = "=========="
	dateFmt       = "Monday, 2 January 2006 15:04:05"
	altDateFmt    = "Monday, 2 January 06 15:04:05"
	unknownAuthor = "Unknown Author"
)

type Range struct {
	From, To int
}

type Clipping struct {
	Title    string
	Author   string
	Type     string
	Page     *int   `json:",omitempty"`
	Location *Range `json:",omitempty"`
	Datetime time.Time
	Content  string
}

type Reader struct {
	*bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

func (r *Reader) Read() (clipping Clipping, err error) {
	// get title and author
	line, err := r.ReadString('\n')
	if err != nil {
		return
	}

	// skip random bom
	if line[:3] == "\ufeff" {
		line = line[3:]
	}

	lastP := strings.LastIndex(line, "(")
	if lastP > 0 {
		clipping.Title = strings.TrimSpace(line[:lastP-1])
		clipping.Author = flipName(strings.TrimSpace(line[lastP+1 : len(line)-3]))
	} else {
		clipping.Title = strings.TrimSpace(line)
		clipping.Author = unknownAuthor
	}

	// get type, pg, loc, datetime
	line, err = r.ReadString('\n')
	if err != nil {
		return
	}

	// skip "- Your"
	line = line[7:]

	// get type
	parts := strings.SplitN(line, " ", 2)
	clipping.Type = parts[0]
	line = parts[1]

	for _, part := range strings.Split(line, "|") {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "Location") || strings.HasPrefix(part, "location") {
			clipping.Location, _ = parseLocation(part[9:])

		} else if strings.HasPrefix(part, "at location") {
			clipping.Location, _ = parseLocation(part[12:])

		} else if strings.HasPrefix(part, "on Page") || strings.HasPrefix(part, "on page") {
			pg, _ := strconv.Atoi(part[8:])
			clipping.Page = &pg

		} else if strings.HasPrefix(part, "Added on") {
			clipping.Datetime, _ = parseDatetime(part[9:])
		}
	}

	// blank line
	if _, err = r.ReadString('\n'); err != nil {
		return
	}

	for {
		line, err = r.ReadString('\n')
		if err != nil {
			return
		}

		if line[:len(line)-2] == delim {
			break
		}

		clipping.Content += line
	}

	clipping.Content = strings.TrimSpace(clipping.Content)

	return
}

func (r *Reader) ReadAll() (clippings []Clipping, err error) {
	for {
		clipping, err := r.Read()
		if err == io.EOF {
			return clippings, nil
		}
		if err != nil {
			return clippings, err
		}
		clippings = append(clippings, clipping)
	}
}

func parseLocation(s string) (*Range, error) {
	fields := strings.Split(s, "-")
	switch len(fields) {
	case 1:
		loc, _ := strconv.Atoi(fields[0])
		return &Range{loc, loc}, nil
	case 2:
		from, _ := strconv.Atoi(fields[0])
		to, _ := strconv.Atoi(fields[1])
		return &Range{from, to}, nil
	}

	return nil, fmt.Errorf("Not in range format: %v\n", s)
}

func parseDatetime(s string) (time.Time, error) {
	t, err := time.Parse(dateFmt, s)
	if err != nil {
		return time.Parse(altDateFmt, s)
	}

	return t, nil
}

func flipName(name string) string {
	parts := strings.SplitN(name, ", ", 2)
	if len(parts) < 2 {
		return name
	}

	return parts[1] + " " + parts[0]
}
