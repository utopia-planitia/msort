package msort

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type manifest struct {
	Kind     string
	Metadata struct {
		Name      string
		Namespace string
	}
	Yaml string
}

func NewManifest(yml string) (manifest, error) {
	d := yaml.NewDecoder(strings.NewReader(yml))
	m := manifest{}

	err := d.Decode(&m)
	if err != nil {
		return manifest{}, err
	}

	m.Yaml = yml

	return m, nil
}

func (m manifest) Print() string {
	return m.Yaml
}

func (m *manifest) SortByKeys() error {
	dec := yaml.NewDecoder(strings.NewReader(m.Yaml))

	var doc yaml.Node
	err := dec.Decode(&doc)
	if err != nil {
		return fmt.Errorf("decode yaml: %v", err)
	}

	sortYAML(&doc)

	out := bytes.NewBuffer(nil)

	enc := yaml.NewEncoder(out)
	defer enc.Close()

	enc.SetIndent(2)

	err = enc.Encode(&doc)
	if err != nil {
		return fmt.Errorf("encode yaml: %v", err)
	}

	m.Yaml = out.String()

	return nil
}

type byKey []*yaml.Node

func (i byKey) Len() int { return len(i) / 2 }

func (i byKey) Swap(x, y int) {
	x *= 2
	y *= 2
	i[x], i[y] = i[y], i[x]         // keys
	i[x+1], i[y+1] = i[y+1], i[x+1] // values
}

func (i byKey) Less(x, y int) bool {
	x *= 2
	y *= 2
	return i[x].Value < i[y].Value
}

func sortYAML(node *yaml.Node) *yaml.Node {
	for i, n := range node.Content {
		node.Content[i] = sortYAML(n)
	}

	if node.Kind == yaml.MappingNode {
		sort.Sort(byKey(node.Content))
	}

	return node
}
