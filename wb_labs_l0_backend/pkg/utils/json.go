package utils

import "encoding/json"

func MarshalSafe(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
