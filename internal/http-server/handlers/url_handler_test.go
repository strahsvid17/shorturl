package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidJSON(t *testing.T) {
	for _, body := range []string{
		``,
		`{"full_url":`,
		`{"full_url":"https://example.com",}`,
		`{"full_url":123}`,
		`[]`,
	} {
		t.Run(body, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			response := httptest.NewRecorder()
			NewURLHandler(nil).CreateURLAlias().ServeHTTP(response, request)
			if response.Code != http.StatusNotAcceptable {
				t.Fatalf("status = %d; want %d", response.Code, http.StatusNotAcceptable)
			}
			if response.Body.String() != "invalid json\n" {
				t.Fatalf("body = %q; want invalid json error", response.Body.String())
			}
		})
	}
}

func TestCreateURLAliasInvalidURL(t *testing.T) {
	for _, url := range []string{
		"",
		"example.com",
		"/path",
		"ftp://example.com",
		"https:///path",
		"http://:9090",
		"http://example.com:abc",
		"https://example.com/%zz",
		"https://example.com/a b",
		" https://example.com",
		"https://example.com/\n",
	} {
		t.Run(url, func(t *testing.T) {
			body, err := json.Marshal(aliasCreateBody{FullURL: url})
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			response := httptest.NewRecorder()
			NewURLHandler(nil).CreateURLAlias().ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d; want %d", response.Code, http.StatusBadRequest)
			}
			if response.Body.String() != "invalid URL\n" {
				t.Fatalf("body = %q; want invalid URL error", response.Body.String())
			}
		})
	}
}
