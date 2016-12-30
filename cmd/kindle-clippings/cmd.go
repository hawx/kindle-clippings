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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"hawx.me/code/kindle-clippings/clippings"
)

const clippingsPath = "documents/My Clippings.txt"

const helpMsg = `Usage: kindle-clippings PATH [--only TYPE]

  Reads clippings from your Kindle and outputs them in JSON format to STDOUT.

   PATH
      Path to Kindle, for example /media/johndoe/Kindle or /Volumes/Kindle.

   --only TYPE
      Only list items of the given type (Bookmark, Note or Highlight).
`

func main() {
	onlyType := flag.String("only", "", "")

	flag.Usage = func() { fmt.Println(helpMsg) }
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(helpMsg)
		return
	}

	if err := run(flag.Arg(0), *onlyType); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(path, onlyType string) error {
	file, err := os.Open(filepath.Join(path, clippingsPath))
	if err != nil {
		file, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	r := clippings.NewReader(file)
	items, err := r.ReadAll()
	if err != nil {
		return err
	}

	if onlyType != "" {
		items = clippings.Filter(items, onlyType)
	}

	return json.NewEncoder(os.Stdout).Encode(items)
}
