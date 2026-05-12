package validation

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/manala/manala/internal/validation"
)

func WithLocator(bytes []byte) validation.ValidateOption {
	return validation.WithLocator(Locator{Bytes: bytes})
}

// Locator resolves a JSON pointer (RFC 6901) to a line/column position within a JSON document.
type Locator struct {
	Bytes []byte
}

var pointerReplacer = strings.NewReplacer("~1", "/", "~0", "~")

func (l Locator) At(location string) (int, int) {
	if location == "" {
		return 0, 0
	}

	tokens := strings.Split(strings.TrimPrefix(location, "/"), "/")
	for i, t := range tokens {
		tokens[i] = pointerReplacer.Replace(t)
	}

	dec := json.NewDecoder(bytes.NewReader(l.Bytes))
	if offset, ok := l.walk(dec, tokens); ok {
		chunk := l.Bytes[:offset]
		return 1 + bytes.Count(chunk, []byte("\n")), int(offset) - bytes.LastIndexByte(chunk, '\n')
	}

	return 0, 0
}

func (l Locator) walk(dec *json.Decoder, tokens []string) (int64, bool) {
	if len(tokens) == 0 {
		return int64(len(l.Bytes)) - int64(len(bytes.TrimLeft(l.Bytes[dec.InputOffset():], " \t\n\r:,"))), true
	}
	if tok, err := dec.Token(); err == nil {
		switch tok {
		case json.Delim('{'):
			for dec.More() {
				k, _ := dec.Token()
				if k.(string) == tokens[0] {
					return l.walk(dec, tokens[1:])
				}
				l.skipValue(dec)
			}
		case json.Delim('['):
			idx, _ := strconv.Atoi(tokens[0])
			for i := 0; dec.More(); i++ {
				if i == idx {
					return l.walk(dec, tokens[1:])
				}
				l.skipValue(dec)
			}
		}
	}
	return 0, false
}

func (l Locator) skipValue(dec *json.Decoder) {
	if t, _ := dec.Token(); t != json.Delim('{') && t != json.Delim('[') {
		return
	}
	for depth := 1; depth > 0; {
		if t, _ := dec.Token(); t == json.Delim('{') || t == json.Delim('[') {
			depth++
		} else if t == json.Delim('}') || t == json.Delim(']') {
			depth--
		}
	}
}
