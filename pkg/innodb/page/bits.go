package page

import (
	"fmt"
	"strings"
)

type Bits []byte

func (t Bits) String() string {
	var elems []string

	for _, i := range t {
		elems = append(elems, fmt.Sprintf("%08b", i))
	}

	return strings.Join(elems, "")
}

func (t Bits) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
