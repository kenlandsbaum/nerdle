package server

import (
	"fmt"
	"testing"
)

func Test_marshaToJsonSuccess(t *testing.T) {
	type testAddress struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	type testStruct struct {
		Name    string      `json:"name"`
		Address testAddress `json:"address"`
	}

	testValue := testStruct{Name: "ken", Address: testAddress{Street: "123 elm st", City: "plano"}}
	expected := `{"name":"ken","address":{"street":"123 elm st","city":"plano"}}`

	result, _ := marshalToJson(testValue)

	if string(result) != expected {
		t.Fatalf("expected %s but got %s\n", expected, string(result))
	}
}

func Test_marshalToJsonFail(t *testing.T) {
	res, err := marshalToJson(make(chan int))
	fmt.Println("res:", res)
	if err == nil {
		t.Fatal("expected non nil error")
	}
}
