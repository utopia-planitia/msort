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

		_, err := fmt.Printf("---\n%s\n", manifest.yaml)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
