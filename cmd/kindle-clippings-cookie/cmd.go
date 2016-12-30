// Kindle Clippings Cookie reads a JSON file of clippings and converts them to
// the formats required for use with fortune(6).
//
// Usage:
//
//     kindle-clippings-cookie OUTPATH < INPATH.json
//
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"hawx.me/code/kindle-clippings/clippings"
	"hawx.me/code/kindle-clippings/fortune"
	"hawx.me/code/kindle-clippings/strfile"
)

const helpMsg = `Usage: kindle-clippings-cookie [--width NUM] OUTPATH

  Reads a JSON clippings file from STDIN and writes two files for use with
  fortune(6) named OUTPATH and OUTPATH.dat.

  This will only output clippings of type Highlight.

   --width NUM=80
      Set the max line length, by default 80 characters.
`

func main() {
	width := flag.Int("width", 80, "")

	flag.Usage = func() { fmt.Println(helpMsg) }
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(helpMsg)
		return
	}

	if err := run(flag.Arg(0), *width); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(outpath string, width int) error {
	var items []clippings.Clipping
	if err := json.NewDecoder(os.Stdin).Decode(&items); err != nil {
		return err
	}

	items = clippings.Filter(items, "Highlight")

	fortuneFile, err := openFile(outpath)
	if err != nil {
		return err
	}
	defer fortuneFile.Close()

	datFile, err := openFile(outpath + ".dat")
	if err != nil {
		return err
	}
	defer datFile.Close()

	var buf bytes.Buffer
	if err = (fortune.Fortunes{
		Items: items,
		Width: width,
	}).WriteTo(&buf); err != nil {
		return err
	}

	tee := io.TeeReader(&buf, fortuneFile)
	return strfile.Strfile(tee, datFile)
}

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}
