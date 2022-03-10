package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type helloWorldhandler struct {
	http.Handler
}

func (h helloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if values != nil {
		fmt.Fprintf(w, "r.URL.Query() :\n")
		for k, v := range values {
			fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		fmt.Fprintf(w, "End r.URL.Query()\n\n")
	}

	values = r.Form
	if values != nil {
		fmt.Fprintf(w, "r.Form :\n")
		for k, v := range values {
			fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		fmt.Fprintf(w, "End r.Form\n\n")
	}

	fmt.Fprintf(w, "Headers :\n")
	for i, v := range r.Header {
		fmt.Fprintf(w, "%v => %v\n", i, v)
	}
	fmt.Fprintf(w, "End Headers\n\n")

	values = r.PostForm
	if values != nil {
		fmt.Fprintf(w, "PostForm")
		for k, v := range values {
			fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		fmt.Fprintf(w, "End PostForm\n\n")
	}

	fmt.Fprintf(w, "r.PostForm: %v\n", r.PostForm)
	fmt.Fprintf(w, "r.RequestURI: %v\n", r.RequestURI)
	fmt.Fprintf(w, "r.URL.Query(): %v\n", r.URL.Query())
	fmt.Fprintf(w, "r.Form: %v\n", r.Form)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if body != nil {
		fmt.Fprintf(w, "body: %s\n", string(body))
	}

	values, err = url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if values != nil {
		fmt.Fprintf(w, "url.ParseQuery(string(body))")
		for k, v := range values {
			fmt.Fprintf(w, "%v => %v\n", k, v)
		}
		fmt.Fprintf(w, "End url.ParseQuery(string(body))\n\n")
	}

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Proto: %s\n", r.Proto)
	fmt.Fprintf(w, "Remote Addr: %s\n", r.RemoteAddr)
	hn, err := os.Hostname()
	if err == nil {
		fmt.Fprintf(w, "Hostname: %s\n", hn)
	}
}

func main() {
	var h helloWorldhandler
	err := http.ListenAndServe(":8080", h)

	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
