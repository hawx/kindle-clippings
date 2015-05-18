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

	"hawx.me/code/kindle-tools/kindle-clippings/clippings"
)

var (
	onlyType = flag.String("only", "", "")
)

const helpMsg = `Usage: kindle-clippings PATH [--only TYPE]

  Reads clippings from your Kindle and outputs them in json format to Stdout.

 PATH
     Path to Kindle, for example /media/johndoe/Kindle or /Volumes/Kindle.

 --only TYPE
     Only list items of the given type (Bookmark, Note or Highlight)
`

func main() {
	flag.Usage = func() {
		fmt.Println(helpMsg)
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(helpMsg)
		return
	}

	path := filepath.Join(flag.Arg(0), "documents/My Clippings.txt")

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
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

	json.NewEncoder(os.Stdout).Encode(items)
}
