package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTTP_BasicRequest(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	sections := []string{
		"=== REQUEST SUMMARY ===",
		"=== URL INFORMATION ===",
		"=== REQUEST HEADERS ===",
		"=== REQUEST BODY ===",
		"=== FORM DATA ===",
		"=== SERVER INFORMATION ===",
		"=== REQUEST COMPLETED ===",
	}
	for _, section := range sections {
		if !strings.Contains(body, section) {
			t.Errorf("response missing section %q", section)
		}
	}

	if !strings.Contains(body, "Method: GET") {
		t.Error("response missing method")
	}
}

func TestServeHTTP_WithQueryParams(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/path?foo=bar&baz=qux", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "foo = bar") {
		t.Error("response missing query param foo=bar")
	}
	if !strings.Contains(body, "baz = qux") {
		t.Error("response missing query param baz=qux")
	}
}

func TestServeHTTP_WithHeaders(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Custom-Header", "custom-value")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "X-Custom-Header") {
		t.Error("response missing custom header name")
	}
	if !strings.Contains(body, "custom-value") {
		t.Error("response missing custom header value")
	}
}

func TestServeHTTP_WithJSONBody(t *testing.T) {
	handler := helloWorldhandler{}
	jsonBody := `{"key":"value","number":42}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	// JSON should be pretty-printed
	if !strings.Contains(body, `"key": "value"`) {
		t.Error("response missing pretty-printed JSON key")
	}
	if !strings.Contains(body, `"number": 42`) {
		t.Error("response missing pretty-printed JSON number")
	}
}

func TestServeHTTP_WithFormBody(t *testing.T) {
	handler := helloWorldhandler{}
	formBody := "username=alice&password=secret"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formBody))
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "username") {
		t.Error("response missing form field username")
	}
	if !strings.Contains(body, "alice") {
		t.Error("response missing form value alice")
	}
}

func TestServeHTTP_ContentType(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	ct := rec.Result().Header.Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Errorf("expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
}

func TestServeHTTP_ContentLength(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cl := rec.Result().Header.Get("Content-Length")
	if cl == "" {
		t.Error("expected Content-Length header to be set")
	}
	body := rec.Body.String()
	if cl != fmt.Sprintf("%d", len(body)) {
		t.Errorf("Content-Length %s does not match body length %d", cl, len(body))
	}
}

func TestGetRealIP_XForwardedFor(t *testing.T) {
	handler := helloWorldhandler{}

	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"single IP", "203.0.113.1", "203.0.113.1"},
		{"multiple IPs", "203.0.113.1, 70.41.3.18, 150.172.238.178", "203.0.113.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", tt.header)
			req.RemoteAddr = "10.0.0.1:1234"

			got := handler.getRealIP(req)
			if got != tt.expected {
				t.Errorf("getRealIP() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetRealIP_XRealIP(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "198.51.100.1")
	req.RemoteAddr = "10.0.0.1:1234"

	got := handler.getRealIP(req)
	if got != "198.51.100.1" {
		t.Errorf("getRealIP() = %q, want %q", got, "198.51.100.1")
	}
}

func TestGetRealIP_CfConnectingIP(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cf-Connecting-Ip", "172.16.0.1")
	req.RemoteAddr = "10.0.0.1:1234"

	got := handler.getRealIP(req)
	if got != "172.16.0.1" {
		t.Errorf("getRealIP() = %q, want %q", got, "172.16.0.1")
	}
}

func TestGetRealIP_Fallback(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:5678"

	got := handler.getRealIP(req)
	if got != "192.168.1.1:5678" {
		t.Errorf("getRealIP() = %q, want %q", got, "192.168.1.1:5678")
	}
}

func BenchmarkServeHTTP(b *testing.B) {
	handler := helloWorldhandler{}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		req := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("X-Custom", "bench")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
