package util

import "encoding/json"

// MergeJSONData merges two JSON objects, with keys in override taking precedence
// over keys in base. Both inputs must be JSON objects ({...}).
func MergeJSONData(base, override []byte) ([]byte, error) {
	var baseMap, overrideMap map[string]any
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(override, &overrideMap); err != nil {
		return nil, err
	}
	for k, v := range overrideMap {
		baseMap[k] = v
	}
	return json.Marshal(baseMap)
}
