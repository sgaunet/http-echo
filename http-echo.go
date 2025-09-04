// Package main provides a simple HTTP echo server that prints request information.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

type requestInfo struct {
	StartTime    time.Time
	Method       string
	URL          string
	Host         string
	RemoteAddr   string
	RealIP       string
	UserAgent    string
	ContentType  string
	ContentLength int64
	Headers      http.Header
	Body         []byte
	QueryParams  url.Values
	FormData     url.Values
	PostFormData url.Values
}

type helloWorldhandler struct {
	http.Handler
}

func (h helloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Collect request information
	info := h.collectRequestInfo(r, startTime)
	
	// Set response content type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	// Print structured output
	h.printRequestSummary(w, info)
	h.printURLInfo(w, info)
	h.printHeaders(w, info)
	h.printRequestBody(w, info)
	h.printFormData(w, info)
	h.printServerInfo(w, info, startTime)
	
	_, _ = fmt.Fprintf(w, "\n=== REQUEST COMPLETED ===")
	_, _ = fmt.Fprintf(w, "\nProcessing Time: %v\n", time.Since(startTime))
}

func (h helloWorldhandler) collectRequestInfo(r *http.Request, startTime time.Time) requestInfo {
	// Parse form data
	_ = r.ParseForm()
	
	// Read body
	body, _ := io.ReadAll(r.Body)
	
	// Get real IP
	realIP := h.getRealIP(r)
	
	return requestInfo{
		StartTime:     startTime,
		Method:        r.Method,
		URL:           r.URL.String(),
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		RealIP:        realIP,
		UserAgent:     r.UserAgent(),
		ContentType:   r.Header.Get("Content-Type"),
		ContentLength: r.ContentLength,
		Headers:       r.Header,
		Body:          body,
		QueryParams:   r.URL.Query(),
		FormData:      r.Form,
		PostFormData:  r.PostForm,
	}
}

func (h helloWorldhandler) getRealIP(r *http.Request) string {
	// Check common reverse proxy headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if parts := strings.Split(ip, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("Cf-Connecting-Ip"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func (h helloWorldhandler) printRequestSummary(w http.ResponseWriter, info requestInfo) {
	_, _ = fmt.Fprintf(w, "=== REQUEST SUMMARY ===\n")
	_, _ = fmt.Fprintf(w, "Timestamp: %s\n", info.StartTime.UTC().Format(time.RFC3339))
	_, _ = fmt.Fprintf(w, "Method: %s | Protocol: HTTP/1.1 | Host: %s\n", info.Method, info.Host)
	if info.URL != "" {
		_, _ = fmt.Fprintf(w, "Full URL: %s\n", info.URL)
	}
	_, _ = fmt.Fprintf(w, "Remote Address: %s\n", info.RemoteAddr)
	if info.RealIP != info.RemoteAddr {
		_, _ = fmt.Fprintf(w, "Real Client IP: %s\n", info.RealIP)
	}
	if info.UserAgent != "" {
		_, _ = fmt.Fprintf(w, "User Agent: %s\n", info.UserAgent)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) printURLInfo(w http.ResponseWriter, info requestInfo) {
	_, _ = fmt.Fprintf(w, "=== URL INFORMATION ===\n")
	_, _ = fmt.Fprintf(w, "Path: %s\n", info.URL)
	
	if len(info.QueryParams) > 0 {
		_, _ = fmt.Fprintf(w, "Query Parameters:\n")
		for key, values := range info.QueryParams {
			for _, value := range values {
				_, _ = fmt.Fprintf(w, "  %s = %s\n", key, value)
			}
		}
	} else {
		_, _ = fmt.Fprintf(w, "Query Parameters: (none)\n")
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) printHeaders(w http.ResponseWriter, info requestInfo) {
	_, _ = fmt.Fprintf(w, "=== REQUEST HEADERS ===\n")
	
	// Highlight important headers first
	importantHeaders := []string{"Content-Type", "Content-Length", "Authorization", "Accept", "Accept-Encoding"}
	for _, headerName := range importantHeaders {
		if value := info.Headers.Get(headerName); value != "" {
			_, _ = fmt.Fprintf(w, "* %-15s: %s\n", headerName, value)
		}
	}
	
	_, _ = fmt.Fprintf(w, "\nAll Headers:\n")
	for name, values := range info.Headers {
		for _, value := range values {
			_, _ = fmt.Fprintf(w, "  %-20s: %s\n", name, value)
		}
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) printRequestBody(w http.ResponseWriter, info requestInfo) {
	_, _ = fmt.Fprintf(w, "=== REQUEST BODY ===\n")
	_, _ = fmt.Fprintf(w, "Content-Length: %d bytes\n", info.ContentLength)
	_, _ = fmt.Fprintf(w, "Content-Type: %s\n", info.ContentType)
	
	if len(info.Body) > 0 {
		_, _ = fmt.Fprintf(w, "Body Content:\n")
		h.formatBodyContent(w, info)
		h.parseBodyAsForm(w, info)
	} else {
		_, _ = fmt.Fprintf(w, "Body: (empty)\n")
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) printFormData(w http.ResponseWriter, info requestInfo) {
	_, _ = fmt.Fprintf(w, "=== FORM DATA ===\n")
	
	if len(info.FormData) > 0 {
		_, _ = fmt.Fprintf(w, "Combined Form Data (GET + POST):\n")
		for key, values := range info.FormData {
			for _, value := range values {
				_, _ = fmt.Fprintf(w, "  %s = %s\n", key, value)
			}
		}
	} else {
		_, _ = fmt.Fprintf(w, "Form Data: (none)\n")
	}
	
	if len(info.PostFormData) > 0 {
		_, _ = fmt.Fprintf(w, "\nPOST Form Data Only:\n")
		for key, values := range info.PostFormData {
			for _, value := range values {
				_, _ = fmt.Fprintf(w, "  %s = %s\n", key, value)
			}
		}
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) printServerInfo(w http.ResponseWriter, _ requestInfo, startTime time.Time) {
	_, _ = fmt.Fprintf(w, "=== SERVER INFORMATION ===\n")
	
	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		_, _ = fmt.Fprintf(w, "Server Hostname: %s\n", hostname)
	}
	
	// Runtime info
	_, _ = fmt.Fprintf(w, "Go Version: %s\n", runtime.Version())
	_, _ = fmt.Fprintf(w, "Server OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	
	// Environment variables (container/k8s info)
	envVars := []string{"HOSTNAME", "POD_NAME", "POD_NAMESPACE", "CONTAINER_NAME"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			_, _ = fmt.Fprintf(w, "%s: %s\n", envVar, value)
		}
	}
	
	_, _ = fmt.Fprintf(w, "Request Start Time: %s\n", startTime.Format(time.RFC3339Nano))
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) formatBodyContent(w http.ResponseWriter, info requestInfo) {
	// Try to format JSON if it looks like JSON
	if strings.Contains(strings.ToLower(info.ContentType), "json") {
		var jsonData interface{}
		if err := json.Unmarshal(info.Body, &jsonData); err == nil {
			if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				_, _ = fmt.Fprintf(w, "%s\n", string(formatted))
				return
			}
		}
	}
	_, _ = fmt.Fprintf(w, "%s\n", string(info.Body))
}

func (h helloWorldhandler) parseBodyAsForm(w http.ResponseWriter, info requestInfo) {
	// Try to parse as form data if it looks like form data
	if parsedValues, err := url.ParseQuery(string(info.Body)); err == nil && len(parsedValues) > 0 {
		_, _ = fmt.Fprintf(w, "\nParsed as form data:\n")
		for key, values := range parsedValues {
			for _, value := range values {
				_, _ = fmt.Fprintf(w, "  %s = %s\n", key, value)
			}
		}
	}
}

const (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 60 * time.Second
	maxHeaderBytes  = 1 << 20
	shutdownTimeout = 5 * time.Second
)

func main() {
	var h helloWorldhandler
	server := &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %s\n", err.Error())
	}
}
