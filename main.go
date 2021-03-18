package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	yml, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err.Error())
	}

	manifests := splitManifests(string(yml))

	sort.Sort(ByOrder(manifests))

	for _, manifest := range manifests {
		if os.Getenv("KEEP_TESTS") == "" {
			if strings.Contains(strings.ToLower(manifest.Metadata.Name), "test") {
				continue
			}
		}

		yaml, err := manifest.Print()
		if err != nil {
			log.Fatalln(err.Error())
		}

		if os.Getenv("DISABLE_KEY_SORTING") == "" {
			yaml, err = manifest.SortedByKeys()
			if err != nil {
				log.Fatalln(err.Error())
			}
		}

		_, err = fmt.Printf("---\n%s\n", yaml)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
