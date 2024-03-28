package server

import (
	"encoding/json"
	"io"
	"log"
)

func marshalToJson(v any) ([]byte, error) {
	bts, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func mustMarshal(v any) []byte {
	bts, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return bts
}

func unmarshalToType[T any](bts []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(bts, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func decodeToType[T any](body io.Reader) (*T, error) {
	var v T
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseRequestBody[T any](body io.ReadCloser) (*T, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	parsed, err := unmarshalToType[T](bodyBytes)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func decodeRequestBody[T any](body io.ReadCloser) (*T, error) {
	decoded, err := decodeToType[T](body)
	if err != nil {
		return nil, err
	}
	return decoded, err
}
