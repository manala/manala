package clean

import (
	"fmt"
)

func YamlStringMap(in map[string]interface{}) map[string]interface{} {
	for k, v := range in {
		in[k] = YamlValue(v)
	}
	return in
}

func YamlValue(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		for i, v := range in {
			in[i] = YamlValue(v)
		}
		return in
	case map[string]interface{}:
		return YamlStringMap(in)
	case map[interface{}]interface{}:
		out := make(map[string]interface{})
		for k, v := range in {
			out[fmt.Sprintf("%v", k)] = YamlValue(v)
		}
		return out
	default:
		return in
	}
}
