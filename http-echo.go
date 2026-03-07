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
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var version = "dev"

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

	// Limit request body size to prevent memory exhaustion from large bodies or forms
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// Collect request information
	info := h.collectRequestInfo(r, startTime)

	// Build response in memory to reduce allocations
	var sb strings.Builder
	sb.Grow(initialBufferSize)

	// Print structured output
	h.printRequestSummary(&sb, info)
	h.printURLInfo(&sb, info)
	h.printHeaders(&sb, info)
	h.printRequestBody(&sb, info)
	h.printFormData(&sb, info)
	h.printServerInfo(&sb, info, startTime)

	_, _ = fmt.Fprintf(&sb, "\n=== REQUEST COMPLETED ===")
	_, _ = fmt.Fprintf(&sb, "\nProcessing Time: %v\n", time.Since(startTime))

	// Write complete response in a single call
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(sb.Len()))
	if _, err := io.WriteString(w, sb.String()); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func (h helloWorldhandler) collectRequestInfo(r *http.Request, startTime time.Time) requestInfo {
	defer func() { _ = r.Body.Close() }()

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing form data: %v", err)
	}
	
	// Read body (size already limited by MaxBytesReader in ServeHTTP)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		body = []byte("(body read failed)")
	}
	
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

func (h helloWorldhandler) printRequestSummary(w io.Writer, info requestInfo) {
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

func (h helloWorldhandler) printURLInfo(w io.Writer, info requestInfo) {
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

func (h helloWorldhandler) printHeaders(w io.Writer, info requestInfo) {
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

func (h helloWorldhandler) printRequestBody(w io.Writer, info requestInfo) {
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

func (h helloWorldhandler) printFormData(w io.Writer, info requestInfo) {
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

func (h helloWorldhandler) printServerInfo(w io.Writer, _ requestInfo, startTime time.Time) {
	_, _ = fmt.Fprintf(w, "=== SERVER INFORMATION ===\n")
	
	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		_, _ = fmt.Fprintf(w, "Server Hostname: %s\n", hostname)
	}
	
	// Runtime info
	_, _ = fmt.Fprintf(w, "Server Version: %s\n", version)
	_, _ = fmt.Fprintf(w, "Go Version: %s\n", runtime.Version())
	_, _ = fmt.Fprintf(w, "Server OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	
	// Environment variables (container/k8s info)
	envVars := []string{"HOSTNAME", "POD_NAME", "POD_NAMESPACE"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			_, _ = fmt.Fprintf(w, "%s: %s\n", envVar, value)
		}
	}
	
	_, _ = fmt.Fprintf(w, "Request Start Time: %s\n", startTime.Format(time.RFC3339Nano))
	_, _ = fmt.Fprintf(w, "\n")
}

func (h helloWorldhandler) formatBodyContent(w io.Writer, info requestInfo) {
	// Try to format JSON if it looks like JSON
	if strings.Contains(strings.ToLower(info.ContentType), "json") {
		var jsonData any
		if err := json.Unmarshal(info.Body, &jsonData); err == nil {
			if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				_, _ = fmt.Fprintf(w, "%s\n", string(formatted))
				return
			}
		}
	}
	_, _ = fmt.Fprintf(w, "%s\n", string(info.Body))
}

func (h helloWorldhandler) parseBodyAsForm(w io.Writer, info requestInfo) {
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
	initialBufferSize  = 2048
	readTimeout        = 10 * time.Second
	readHeaderTimeout  = 5 * time.Second
	writeTimeout       = 10 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20
	maxBodySize       = 10 << 20 // 10MB
	shutdownTimeout   = 5 * time.Second
)

func newServer() *http.Server {
	var h helloWorldhandler
	return &http.Server{
		Addr:              ":8080",
		Handler:           h,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func main() {
	server := newServer()

	go func() {
		log.Printf("Starting server on :8080 (version %s)", version)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received, shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %s\n", err.Error())
	}
	log.Println("Server stopped")
}
