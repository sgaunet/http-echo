package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestNewServer_Timeouts(t *testing.T) {
	srv := newServer()

	tests := []struct {
		name string
		got  time.Duration
		want time.Duration
	}{
		{"ReadTimeout", srv.ReadTimeout, readTimeout},
		{"ReadHeaderTimeout", srv.ReadHeaderTimeout, readHeaderTimeout},
		{"WriteTimeout", srv.WriteTimeout, writeTimeout},
		{"IdleTimeout", srv.IdleTimeout, idleTimeout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestNewServer_ReadHeaderTimeoutLessThanReadTimeout(t *testing.T) {
	srv := newServer()
	if srv.ReadHeaderTimeout >= srv.ReadTimeout {
		t.Errorf("ReadHeaderTimeout (%v) should be less than ReadTimeout (%v)",
			srv.ReadHeaderTimeout, srv.ReadTimeout)
	}
}

func TestNewServer_Addr(t *testing.T) {
	srv := newServer()
	if srv.Addr != ":8080" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":8080")
	}
}

func TestNewServer_MaxHeaderBytes(t *testing.T) {
	srv := newServer()
	if srv.MaxHeaderBytes != maxHeaderBytes {
		t.Errorf("MaxHeaderBytes = %d, want %d", srv.MaxHeaderBytes, maxHeaderBytes)
	}
}

func TestGetRealIP_Priority(t *testing.T) {
	handler := helloWorldhandler{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.1.1.1")
	req.Header.Set("X-Real-IP", "2.2.2.2")
	req.Header.Set("Cf-Connecting-Ip", "3.3.3.3")
	req.RemoteAddr = "4.4.4.4:1234"

	got := handler.getRealIP(req)
	if got != "1.1.1.1" {
		t.Errorf("X-Forwarded-For should take priority, got %q", got)
	}
}

func TestFormatBodyContent_JSON(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		ContentType: "application/json",
		Body:        []byte(`{"b":2,"a":1}`),
	}

	handler.formatBodyContent(&sb, info)
	out := sb.String()

	if !strings.Contains(out, `"a": 1`) || !strings.Contains(out, `"b": 2`) {
		t.Errorf("expected pretty-printed JSON, got %q", out)
	}
}

func TestFormatBodyContent_InvalidJSON(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		ContentType: "application/json",
		Body:        []byte(`{invalid`),
	}

	handler.formatBodyContent(&sb, info)

	if !strings.Contains(sb.String(), "{invalid") {
		t.Error("invalid JSON should be output as-is")
	}
}

func TestFormatBodyContent_PlainText(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		ContentType: "text/plain",
		Body:        []byte("hello world"),
	}

	handler.formatBodyContent(&sb, info)

	if !strings.Contains(sb.String(), "hello world") {
		t.Error("expected plain text in output")
	}
}

func TestParseBodyAsForm(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		Body: []byte("key1=val1&key2=val2"),
	}

	handler.parseBodyAsForm(&sb, info)
	out := sb.String()

	if !strings.Contains(out, "Parsed as form data:") {
		t.Error("expected form data parsing header")
	}
	if !strings.Contains(out, "key1 = val1") {
		t.Error("missing key1 in parsed form data")
	}
	if !strings.Contains(out, "key2 = val2") {
		t.Error("missing key2 in parsed form data")
	}
}

func TestParseBodyAsForm_EmptyBody(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		Body: []byte(""),
	}

	handler.parseBodyAsForm(&sb, info)

	if sb.Len() != 0 {
		t.Errorf("empty body should produce no output, got %q", sb.String())
	}
}

func TestPrintRequestSummary_RealIPShown(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		Method:     "GET",
		Host:       "localhost:8080",
		URL:        "/test",
		RemoteAddr: "127.0.0.1:1234",
		RealIP:     "10.0.0.1",
		UserAgent:  "test-agent",
	}

	handler.printRequestSummary(&sb, info)
	out := sb.String()

	if !strings.Contains(out, "Real Client IP: 10.0.0.1") {
		t.Error("should show real IP when different from remote addr")
	}
	if !strings.Contains(out, "User Agent: test-agent") {
		t.Error("missing user agent")
	}
}

func TestPrintRequestSummary_RealIPHidden(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		RemoteAddr: "127.0.0.1:1234",
		RealIP:     "127.0.0.1:1234",
	}

	handler.printRequestSummary(&sb, info)

	if strings.Contains(sb.String(), "Real Client IP") {
		t.Error("should not show real IP when same as remote addr")
	}
}

func TestPrintURLInfo_NoQueryParams(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{URL: "/simple"}

	handler.printURLInfo(&sb, info)

	if !strings.Contains(sb.String(), "Query Parameters: (none)") {
		t.Error("should show (none) for empty query params")
	}
}

func TestPrintHeaders_ImportantHighlighted(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer token")
	headers.Set("X-Custom", "value")
	info := requestInfo{Headers: headers}

	handler.printHeaders(&sb, info)
	out := sb.String()

	if !strings.Contains(out, "* Content-Type") {
		t.Error("Content-Type should be highlighted")
	}
	if !strings.Contains(out, "* Authorization") {
		t.Error("Authorization should be highlighted")
	}
	if !strings.Contains(out, "X-Custom") {
		t.Error("missing custom header")
	}
}

func TestPrintRequestBody_Empty(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{}

	handler.printRequestBody(&sb, info)

	if !strings.Contains(sb.String(), "Body: (empty)") {
		t.Error("should show empty body message")
	}
}

func TestPrintFormData_Empty(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{}

	handler.printFormData(&sb, info)

	if !strings.Contains(sb.String(), "Form Data: (none)") {
		t.Error("should show (none) for empty form data")
	}
}

func TestPrintFormData_WithPostData(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder
	info := requestInfo{
		FormData:     map[string][]string{"field": {"val"}},
		PostFormData: map[string][]string{"field": {"val"}},
	}

	handler.printFormData(&sb, info)
	out := sb.String()

	if !strings.Contains(out, "Combined Form Data (GET + POST):") {
		t.Error("missing combined form data section")
	}
	if !strings.Contains(out, "POST Form Data Only:") {
		t.Error("missing POST form data section")
	}
}

func TestPrintServerInfo(t *testing.T) {
	handler := helloWorldhandler{}
	var sb strings.Builder

	handler.printServerInfo(&sb, requestInfo{}, time.Now())
	out := sb.String()

	if !strings.Contains(out, "=== SERVER INFORMATION ===") {
		t.Error("missing section header")
	}
	if !strings.Contains(out, "Go Version:") {
		t.Error("missing Go version")
	}
	if !strings.Contains(out, "Server OS:") {
		t.Error("missing server OS")
	}
}

func TestServeHTTP_BodySizeLimit(t *testing.T) {
	handler := helloWorldhandler{}
	// Create a body larger than maxBodySize (10MB)
	largeBody := strings.NewReader(strings.Repeat("x", 11<<20))
	req := httptest.NewRequest(http.MethodPost, "/", largeBody)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	// The response should still complete (not panic), body should be truncated or show error
	if !strings.Contains(body, "=== REQUEST COMPLETED ===") {
		t.Error("response should complete even with oversized body")
	}
}

func BenchmarkServeHTTP(b *testing.B) {
	handler := helloWorldhandler{}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		req := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		req.Header.Set("X-Custom", "bench")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
