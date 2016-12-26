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
	"log"
	"os"
	"path/filepath"

	"hawx.me/code/kindle-clippings/clippings"
	"hawx.me/code/kindle-clippings/fortune"
)

const clippingsPath = "documents/My Clippings.txt"

const helpMsg = `Usage: kindle-clippings PATH [--only TYPE] [--fortune]

  Reads clippings from your Kindle and outputs them in JSON format to STDOUT.

   PATH
      Path to Kindle, for example /media/johndoe/Kindle or /Volumes/Kindle.

   --only TYPE
      Only list items of the given type (Bookmark, Note or Highlight).

   --fortune
      Output in a format for use as a fortune(6) cookie file.

      After creating the file you will need to run 'strfile' over it to produce
      a .dat companion.
`

func main() {
	var (
		onlyType      = flag.String("only", "", "")
		fortuneOutput = flag.Bool("fortune", false, "")
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
		fmt.Println("err: %v", err)
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

	if *fortuneOutput {
		fortune.Fortunes(items).WriteTo(os.Stdout)
	} else {
		json.NewEncoder(os.Stdout).Encode(items)
	}
}
