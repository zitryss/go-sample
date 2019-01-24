package testing

import (
	"net/http/httptest"
	"testing"
)

func CheckStatusCode(t *testing.T, w *httptest.ResponseRecorder, code int) {
	t.Helper()
	got := w.Code
	want := code
	if got != want {
		t.Errorf("Status Code = %v; want %v", got, want)
	}
}

func CheckContentType(t *testing.T, w *httptest.ResponseRecorder, content string) {
	t.Helper()
	got := w.Result().Header.Get("Content-Type")
	want := content
	if got != want {
		t.Errorf("Content-Type = %v; want %v", got, want)
	}
}

func CheckBody(t *testing.T, w *httptest.ResponseRecorder, body string) {
	got := w.Body.String()
	want := body
	if got != want {
		t.Errorf("Body = %v; want %v", got, want)
	}
}
