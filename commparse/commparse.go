package commparse

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

// CommitteeID is the unique ID for a committee in the FEC datasets.
type CommitteeID string

type Committee struct {
	ID   CommitteeID
	Name string
}

// ParseFile parses a file of committees and their info into a slice of
// Committee structs.
func ParseFile(fp string) ([]Committee, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(b, []byte{'\n'})
	pacs := make([]Committee, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		cols := bytes.Split(line, []byte{'|'})
		if len(cols) != 15 {
			return nil, fmt.Errorf("on line %d, expected 15 columns, got %d", i+1, cols)
		}
		pacs = append(pacs, Committee{
			ID:   CommitteeID(string(cols[0])),
			Name: string(cols[1]),
		})
	}
	return pacs, nil
}
