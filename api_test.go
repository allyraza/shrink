package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_WelcomeHandler(t *testing.T) {
	m := New()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	m.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, Got %v", http.StatusOK, w.Code)
	}

	got := w.Header().Get("Content-Type")
	expected := "text/html"
	if got != expected {
		t.Errorf("Expected %v, Got %v", expected, got)
	}
}

func Test_ShortHandler(t *testing.T) {
	m := New()
	w := httptest.NewRecorder()
	// @todo: send params as post body
	data := url.Values{"url": {"foo"}}
	r, _ := http.NewRequest("POST", "/short", nil)
	r.URL.RawQuery = data.Encode()

	m.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, Got %v", http.StatusOK, w.Code)
	}
}
