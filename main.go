package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Project struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tags        []string `json:"tags"`
	URL         string `json:"url"`
	Repo        string `json:"repo"`
	Featured    bool   `json:"featured"`
}

func loadProjects() ([]Project, error) {
	data, err := fs.ReadFile(contentFS, "content/projects.json")
	if err != nil {
		return nil, fmt.Errorf("read projects.json: %w", err)
	}
	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("unmarshal projects.json: %w", err)
	}
	return projects, nil
}

func getTemplateFS() fs.FS {
	sub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		log.Fatalf("sub templates: %v", err)
	}
	return sub
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmplFS := getTemplateFS()
	tmpl := template.New("base.html")
	var err error

	if name == "index.html" {
		tmpl, err = tmpl.ParseFS(tmplFS, "base.html", name)
	} else {
		tmpl, err = tmpl.ParseFS(tmplFS, "base.html", "contentPane.html", name)
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

type PageData struct {
	Title    string
	Projects []Project
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := loadProjects()
	if err != nil {
		log.Printf("load projects: %v", err)
	}
	renderTemplate(w, "index.html", PageData{
		Title:    "Home",
		Projects: projects,
	})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := loadProjects()
	if err != nil {
		log.Printf("load projects: %v", err)
	}
	renderTemplate(w, "projects.html", PageData{
		Title:    "Projects",
		Projects: projects,
	})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", PageData{
		Title: "Contact",
	})
}

func main() {
	mux := http.NewServeMux()

	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	imagesSub, _ := fs.Sub(imageFS, "images")
	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.FS(imagesSub))))

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /projects", projectsHandler)
	mux.HandleFunc("GET /contact", contactHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
