package server

import "encoding/json"

func marshalToJson(v any) ([]byte, error) {
	bts, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func unmarshalToType[T any](bts []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(bts, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
