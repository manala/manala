package cleaner

import (
	"fmt"
)

func Clean(in map[string]interface{}) map[string]interface{} {
	for k, v := range in {
		in[k] = cleanValue(v)
	}
	return in
}

func cleanValue(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		for i, v := range in {
			in[i] = cleanValue(v)
		}
		return in
	case map[string]interface{}:
		return Clean(in)
	case map[interface{}]interface{}:
		out := make(map[string]interface{})
		for k, v := range in {
			out[fmt.Sprintf("%v", k)] = cleanValue(v)
		}
		return out
	default:
		return in
	}
}
