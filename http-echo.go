// Package main provides a simple HTTP echo server that prints request information.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type helloWorldhandler struct {
	http.Handler
}

func (h helloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.printQueryParams(w, r)
	h.printFormData(w, r)
	h.printHeaders(w, r)
	h.printPostFormData(w, r)
	h.printRequestInfo(w, r)
	h.printBodyAndParsedData(w, r)
	h.printConnectionInfo(w, r)
}

func (h helloWorldhandler) printQueryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values) > 0 {
		_, _ = fmt.Fprintf(w, "r.URL.Query() :\n")
		for k, v := range values {
			_, _ = fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		_, _ = fmt.Fprintf(w, "End r.URL.Query()\n\n")
	}
}

func (h helloWorldhandler) printFormData(w http.ResponseWriter, r *http.Request) {
	values := r.Form
	if len(values) > 0 {
		_, _ = fmt.Fprintf(w, "r.Form :\n")
		for k, v := range values {
			_, _ = fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		_, _ = fmt.Fprintf(w, "End r.Form\n\n")
	}
}

func (h helloWorldhandler) printHeaders(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Headers :\n")
	for i, v := range r.Header {
		_, _ = fmt.Fprintf(w, "%v => %v\n", i, v)
	}
	_, _ = fmt.Fprintf(w, "End Headers\n\n")
}

func (h helloWorldhandler) printPostFormData(w http.ResponseWriter, r *http.Request) {
	values := r.PostForm
	if len(values) > 0 {
		_, _ = fmt.Fprintf(w, "PostForm")
		for k, v := range values {
			_, _ = fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		_, _ = fmt.Fprintf(w, "End PostForm\n\n")
	}
}

func (h helloWorldhandler) printRequestInfo(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "r.PostForm: %v\n", r.PostForm)
	_, _ = fmt.Fprintf(w, "r.RequestURI: %v\n", r.RequestURI)
	_, _ = fmt.Fprintf(w, "r.URL.Query(): %v\n", r.URL.Query())
	_, _ = fmt.Fprintf(w, "r.Form: %v\n", r.Form)
}

func (h helloWorldhandler) printBodyAndParsedData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(body) > 0 {
		_, _ = fmt.Fprintf(w, "body: %s\n", string(body))
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(values) > 0 {
		_, _ = fmt.Fprintf(w, "url.ParseQuery(string(body))")
		for k, v := range values {
			_, _ = fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		_, _ = fmt.Fprintf(w, "End url.ParseQuery(string(body))\n\n")
	}
}

func (h helloWorldhandler) printConnectionInfo(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Method: %s\n", r.Method)
	_, _ = fmt.Fprintf(w, "Host: %s\n", r.Host)
	_, _ = fmt.Fprintf(w, "Proto: %s\n", r.Proto)
	_, _ = fmt.Fprintf(w, "Remote Addr: %s\n", r.RemoteAddr)
	hn, err := os.Hostname()
	if err == nil {
		_, _ = fmt.Fprintf(w, "Hostname: %s\n", hn)
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
