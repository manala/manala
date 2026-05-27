package validation

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/manala/manala/internal/validation"

	"github.com/go-openapi/jsonpointer"
)

func WithLocator(bytes []byte) validation.ValidateOption {
	return validation.WithLocator(Locator{Bytes: bytes})
}

// Locator resolves a JSON pointer (RFC 6901) to a line/column position within a JSON document.
type Locator struct {
	Bytes []byte
}

func (l Locator) ValueAt(location string) (int, int) {
	return l.at(location, false)
}

func (l Locator) PropertyAt(location string) (int, int) {
	return l.at(location, true)
}

// at resolves a JSON pointer and returns the line/column of the matched
// entry's key (when asProperty is true) or value.
func (l Locator) at(location string, asProperty bool) (int, int) {
	p, err := jsonpointer.New(location)
	if err != nil {
		return 0, 0
	}

	property, value, ok := l.walk(json.NewDecoder(bytes.NewReader(l.Bytes)), p.DecodedTokens(), 0)
	if !ok {
		return 0, 0
	}

	offset := value
	if asProperty {
		offset = property
	}

	chunk := l.Bytes[:offset]
	return 1 + bytes.Count(chunk, []byte("\n")), int(offset) - bytes.LastIndexByte(chunk, '\n')
}

// walk descends through `tokens`, returning the property/value byte offsets
// of the final match. `property` carries the parent-level offset of the
// current entry: the key for object members, the element start for array members.
func (l Locator) walk(dec *json.Decoder, tokens []string, property int64) (int64, int64, bool) {
	if len(tokens) == 0 {
		return property, l.skip(dec.InputOffset()), true
	}

	tok, err := dec.Token()
	if err != nil {
		return 0, 0, false
	}

	switch tok {
	case json.Delim('{'):
		for dec.More() {
			keyOffset := l.skip(dec.InputOffset())
			key, _ := dec.Token()
			if key == tokens[0] {
				return l.walk(dec, tokens[1:], keyOffset)
			}
			_ = dec.Decode(&json.RawMessage{})
		}
	case json.Delim('['):
		idx, _ := strconv.Atoi(tokens[0])
		for i := 0; dec.More(); i++ {
			elemOffset := l.skip(dec.InputOffset())
			if i == idx {
				return l.walk(dec, tokens[1:], elemOffset)
			}
			_ = dec.Decode(&json.RawMessage{})
		}
	}

	return 0, 0, false
}

// skip advances past JSON structural whitespace and separators (`:` and `,`).
func (l Locator) skip(offset int64) int64 {
	return int64(len(l.Bytes)) - int64(len(bytes.TrimLeft(l.Bytes[offset:], " \t\n\r:,")))
}
