package utils

import "encoding/json"

func ToJSON(v interface{}) []byte {
	bt, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bt
}
