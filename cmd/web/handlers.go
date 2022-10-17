package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served: %s\n", r.Host)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"home.gohtml", url+"index.gohtml"))
	if r.Method != http.MethodPost {
		_ = tmpl.ExecuteTemplate(w, "base", nil)
		return
	}
}

func typePaymentHandler(w http.ResponseWriter, r *http.Request) {
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"type_payment.gohtml", url+"index.gohtml"))
	if r.Method != http.MethodPost {
		_ = tmpl.ExecuteTemplate(w, "base", nil)
		return
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is:"
	_, _ = fmt.Fprintf(w, "<h1 align=\"centr\">%s</h1>", Body)
	_, _ = fmt.Fprintf(w, "<h2 align=\"centr\">%s</h2>", t)
	_, _ = fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served time for: %s\n", r.Host)
}
