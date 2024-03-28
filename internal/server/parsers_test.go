package server

import (
	"bytes"
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

func Test_decodeToType(t *testing.T) {
	type subThing struct {
		OtherOne string `json:"otherOne"`
		OtherTwo string `json:"otherTwo"`
	}
	type thing struct {
		Name     string   `json:"name"`
		Age      int      `json:"age"`
		SubThing subThing `json:"subThing"`
	}
	testBody := []byte(`{"name":"ken","age":100,"subThing":{"otherOne":"lol","otherTwo":"umad"}}`)

	res, err := decodeToType[thing](bytes.NewReader(testBody))
	if err != nil {
		t.Errorf("expected nil error but got %s", err.Error())
	}
	if res.Name != "ken" {
		t.Errorf("got %s but expected %s\n", res.Name, "ken")
	}
	if res.Age != 100 {
		t.Errorf("got %d but expected %d\n", res.Age, 100)
	}
	if res.SubThing.OtherOne != "lol" {
		t.Errorf("got %s but expected %s\n", res.SubThing.OtherOne, "lol")
	}
	if res.SubThing.OtherTwo != "umad" {
		t.Errorf("got %s but expected %s\n", res.SubThing.OtherTwo, "umad")
	}
}
