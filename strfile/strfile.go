// Package strfile provides an implementation of strfile(6), given no arguments.
package strfile

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

const delim = '%'

// Strfile should perform the same operation as would running strfile(6) over a
// file with the Readers contents, given no options.
func Strfile(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanDelim)

	var (
		longlen  uint32
		shortlen uint32 = 0xffffffff
		numstr   uint32
		pos      uint32
		offsets  []uint32
	)

	for scanner.Scan() {
		line := scanner.Bytes()

		numstr += 1
		length := uint32(len(line))
		pos += length

		if bytes.HasSuffix(line, []byte{'\n', delim, '\n'}) {
			length -= 2
		}
		if longlen < length {
			longlen = length
		}
		if length < shortlen {
			shortlen = length
		}

		offsets = append(offsets, pos)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	err := binary.Write(w, binary.BigEndian, header{
		Version:  2,
		Numstr:   numstr,
		Longlen:  longlen,
		Shortlen: shortlen,
		Flags:    0,
		Delim:    [8]byte{delim, 0, 0, 0, 0, 0, 0, 0},
	})
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, offsets)
}

type header struct {
	Version  uint32
	Numstr   uint32
	Longlen  uint32
	Shortlen uint32
	Flags    uint32
	Delim    [8]byte
}

func scanDelim(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte{'\n', delim, '\n'}); i >= 0 {
		return i + 3, data[0 : i+3], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
