package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testFn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"key":"value"}`))
}

func testTextFn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do you even go"))
}

func Test_jsonContent(t *testing.T) {
	expectedBody := `{"key":"value"}`
	expectedHeader := "application/json"

	handler := jsonContent(http.HandlerFunc(testFn))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/test", nil)

	handler.ServeHTTP(w, r)

	header := w.Result().Header["Content-Type"]
	if header[0] != expectedHeader {
		t.Errorf("expected %s but got %s\n", expectedHeader, header[0])
	}
	body := w.Result().Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	actual := string(bodyBytes)
	if actual != expectedBody {
		t.Errorf("expected %s but got %s", expectedBody, actual)
	}
}

func Test_useJsonContent(t *testing.T) {
	expectedBody := `{"key":"value"}`
	expectedHeader := "application/json"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/test", nil)

	handlerFunc := useJsonContent(testFn)

	handlerFunc(w, r)

	header := w.Result().Header["Content-Type"]
	if header[0] != expectedHeader {
		t.Errorf("expected %s but got %s\n", expectedHeader, header[0])
	}
	body := w.Result().Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	actual := string(bodyBytes)
	if actual != expectedBody {
		t.Errorf("expected %s but got %s", expectedBody, actual)
	}
}

func Test_useTextContent(t *testing.T) {
	expectedBody := `do you even go`
	expectedHeader := "text/plain"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/test", nil)

	handlerFunc := useTextContent(testTextFn)

	handlerFunc(w, r)

	header := w.Result().Header["Content-Type"]
	if header[0] != expectedHeader {
		t.Errorf("expected %s but got %s\n", expectedHeader, header[0])
	}
	body := w.Result().Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("expected nil error but got %s", err)
	}
	actual := string(bodyBytes)
	if actual != expectedBody {
		t.Errorf("expected %s but got %s", expectedBody, actual)
	}
}
