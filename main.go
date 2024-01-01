package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/layout/", getWithLayout)
	http.HandleFunc("/", get)

	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		panic(err)
	}
}

func getWithLayout(w http.ResponseWriter, r *http.Request) {
	t, err := getPathTemplate(r.URL.Path, "layout.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	t, err := getPathTemplate(r.URL.Path, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPathTemplate(path string, layout string) (*template.Template, error) {
	base := filepath.Base(path) + ".gohtml"

	if layout == "" {
		return template.New(base).ParseFiles("templates/" + base)
	} else {
		return template.New(layout).ParseFiles("templates/"+layout, "templates/"+base)
	}
}
