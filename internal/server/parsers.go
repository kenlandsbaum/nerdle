package server

import "encoding/json"

func marshalToJson(v any) ([]byte, error) {
	bts, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bts, nil
}
