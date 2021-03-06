package msort

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strings"
)

func Sort(in []byte, sortKeys, dropTests bool) (string, error) {
	manifests := NewManifests(string(in))

	if sortKeys {
		err := manifests.SortByKeys()
		if err != nil {
			return "", fmt.Errorf("sort document by key: %v", err)
		}
	}

	if dropTests {
		manifests.DropTest()
	}

	manifests.SortDocuments()

	return manifests.String(), nil
}

type manifests []manifest

func NewManifests(yml string) manifests {
	yml = strings.TrimSpace(yml)

	if strings.HasPrefix(yml, "---\n") {
		yml = "\n" + yml
	}

	if strings.HasSuffix(yml, "\n---") {
		yml += "\n"
	}

	separator := regexp.MustCompile("\n---\n")
	chunks := separator.Split(yml, -1)
	manifests := []manifest{}

	for _, doc := range chunks {
		doc = strings.TrimSpace(doc)

		if doc == "" {
			continue
		}

		manifest, err := NewManifest(doc)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Panicf("parse yaml failed: %v: %v", err, doc)
		}

		manifests = append(manifests, manifest)
	}

	return manifests
}

func (m manifests) SortByKeys() error {
	for i := range m {
		err := m[i].SortByKeys()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manifests) DropTest() {
	filtered := []manifest{}

	for _, manifest := range *m {
		if strings.Contains(strings.ToLower(manifest.Metadata.Name), "test") {
			continue
		}

		filtered = append(filtered, manifest)
	}

	*m = filtered
}

func (m manifests) SortDocuments() {
	sort.Sort(byKind(m))
}

type byKind []manifest

func (a byKind) Len() int      { return len(a) }
func (a byKind) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byKind) Less(i, j int) bool {
	if a[i].Kind != a[j].Kind {
		return a[i].Kind < a[j].Kind
	}

	if a[i].Metadata.Namespace != a[j].Metadata.Namespace {
		return a[i].Metadata.Namespace < a[j].Metadata.Namespace
	}

	if a[i].Metadata.Name != a[j].Metadata.Name {
		return a[i].Metadata.Name < a[j].Metadata.Name
	}

	return a[i].Yaml < a[j].Yaml
}

func (m manifests) String() string {
	buf := strings.Builder{}
	first := true

	for _, manifest := range m {
		yaml := manifest.Print()

		if strings.TrimSpace(yaml) == "" {
			continue
		}

		if !first {
			_, err := buf.WriteString("---\n")
			if err != nil {
				log.Fatalf("write document seperator to buffer: %v", err)
			}
		}

		first = false

		_, err := buf.WriteString(strings.TrimSpace(yaml) + "\n")
		if err != nil {
			log.Fatalf("write document to buffer: %v", err)
		}
	}

	return buf.String()
}
