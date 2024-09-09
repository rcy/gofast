package main

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"simplegohtmx/ws"
)

func main() {
	http.HandleFunc("/", get)
	http.HandleFunc("/subscribe", ws.Subscribe)
	http.HandleFunc("/publish", ws.Publish)

	log.Println("listening on http://localhost:6969")
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	t, err := requestTemplate(r)
	if errors.Is(err, fs.ErrNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
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

func requestTemplate(r *http.Request) (*template.Template, error) {
	base := filepath.Base(r.URL.Path)
	log.Println("base", base)
	if base == "/" {
		base = "index"
	}
	base += ".gohtml"
	dir := "templates" + filepath.Dir(r.URL.Path) + "/"

	if base[0] == '_' {
		return template.New(base).ParseFiles(dir + base)
	} else {
		return template.New("_layout.gohtml").ParseFiles(dir+"_layout.gohtml", dir+base)
	}
}
