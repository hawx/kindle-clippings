// Kindle Clippings reads the clippings file from a Kindle and output them as
// json.
//
// Usage:
//
//     kindle-clippings PATH [--only TYPE]
//
// where PATH is the path to the mounted Kindle. Json output is written to
// Stdout.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"hawx.me/code/kindle-clippings/clippings"
	"hawx.me/code/kindle-clippings/fortune"
	"hawx.me/code/kindle-clippings/strfile"
)

const clippingsPath = "documents/My Clippings.txt"

const helpMsg = `Usage: kindle-clippings PATH [--only TYPE] [--fortune]

  Reads clippings from your Kindle and outputs them in JSON format to STDOUT.

   PATH
      Path to Kindle, for example /media/johndoe/Kindle or /Volumes/Kindle.

   --only TYPE
      Only list items of the given type (Bookmark, Note or Highlight).

   --fortune OUTPATH
      Output in a format for use as a fortune(6) cookie file. This will actually
      create two files: 'OUTPATH' and 'OUTPATH.dat'.
`

func main() {
	var (
		onlyType      = flag.String("only", "", "")
		fortuneOutput = flag.String("fortune", "", "")
	)
	flag.Usage = func() { fmt.Println(helpMsg) }
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(helpMsg)
		return
	}

	file, err := os.Open(filepath.Join(flag.Arg(0), clippingsPath))
	if err != nil {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	r := clippings.NewReader(file)
	items, err := r.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}

	if *onlyType != "" {
		var filtered []clippings.Clipping

		for _, item := range items {
			if item.Type == *onlyType {
				filtered = append(filtered, item)
			}
		}

		items = filtered
	}

	if *fortuneOutput != "" {
		fortuneFile, err := openFile(*fortuneOutput)
		if err != nil {
			log.Println(err)
			return
		}
		defer fortuneFile.Close()

		datFile, err := openFile(*fortuneOutput + ".dat")
		if err != nil {
			log.Println(err)
			return
		}
		defer datFile.Close()

		var buf bytes.Buffer
		fortune.Fortunes(items).WriteTo(&buf)

		tee := io.TeeReader(&buf, fortuneFile)
		strfile.Strfile(tee, datFile)

	} else {
		json.NewEncoder(os.Stdout).Encode(items)
	}
}

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}
