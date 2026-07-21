package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var pages []Page

func getTemplateFS() fs.FS {
	sub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		log.Fatalf("sub templates: %v", err)
	}
	return sub
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}, extraTemplates ...string) {
	tmplFS := getTemplateFS()
	tmpl := template.New("base.html")
	var err error

	if name == "index.html" {
		tmpl, err = tmpl.ParseFS(tmplFS, "base.html", "spiral.html", name)
	} else {
		files := []string{"base.html", "spiral.html", "contentPane.html", name}
		files = append(files, extraTemplates...)
		tmpl, err = tmpl.ParseFS(tmplFS, files...)
	}
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		log.Printf("template parse error: %v", err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func main() {
	pages = initPages()

	mux := http.NewServeMux()

	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	imagesSub, _ := fs.Sub(imageFS, "images")
	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.FS(imagesSub))))

	for _, p := range pages {
		mux.HandleFunc(fmt.Sprintf("GET %v", p.Route), p.Handler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
