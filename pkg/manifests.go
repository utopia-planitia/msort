package msort

import (
	"io"
	"log"
	"regexp"
	"sort"
	"strings"
)

type manifests struct {
	manifests []manifest
}

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
	ls := []manifest{}

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

		ls = append(ls, manifest)
	}

	return manifests{
		manifests: ls,
	}
}

func (m manifests) SortByKeys() error {
	for i := range m.manifests {
		err := m.manifests[i].SortByKeys()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m manifests) DropTest() error {
	filtered := []manifest{}

	for _, manifest := range m.manifests {
		if strings.Contains(strings.ToLower(manifest.Metadata.Name), "test") {
			continue
		}

		filtered = append(filtered, manifest)
	}

	m.manifests = filtered

	return nil
}

func (m manifests) SortDocuments() {
	sort.Sort(byKind(m.manifests))
}

func (m manifests) String() string {
	buf := strings.Builder{}
	first := true

	for _, manifest := range m.manifests {
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
