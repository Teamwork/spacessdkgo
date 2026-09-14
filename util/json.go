package util

import "maps"

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
	maps.Copy(baseMap, overrideMap)
	return json.Marshal(baseMap)
}
