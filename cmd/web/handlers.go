package main

import (
	"fmt"
	"net/http"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served: %s\n", r.Host)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is:"
	_, _ = fmt.Fprintf(w, "<h1 align=\"centr\">%s</h1>", Body)
	_, _ = fmt.Fprintf(w, "<h2 align=\"centr\">%s</h2>", t)
	_, _ = fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served time for: %s\n", r.Host)
}
