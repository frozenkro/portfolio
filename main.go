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

var pages []Page

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
	Repo        string   `json:"repo"`
	Featured    bool     `json:"featured"`
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
		tmpl, err = tmpl.ParseFS(tmplFS, "base.html", "spiral.html", name)
	} else {
		tmpl, err = tmpl.ParseFS(tmplFS, "base.html", "spiral.html", "contentPane.html", name)
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
	Circles  []BrokenCircle
	NavItems []NavItem
}

type BrokenCircle struct {
	Diameter string
	Rotate   int
}

// Desktop specific nav items
// only displayed on 'content' pages
type NavItem struct {
	Active bool
	Last   bool
	Route  string
	Title  string
}

type Page struct {
	Route   string
	Title   string
	Handler func(http.ResponseWriter, *http.Request)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := loadProjects()
	if err != nil {
		log.Printf("load projects: %v", err)
	}
	renderTemplate(w, "index.html", PageData{
		Title:    "Home",
		Projects: projects,
		Circles:  getCircles(),
		NavItems: getNavItems("/"),
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
		Circles:  getCircles(),
		NavItems: getNavItems("/projects"),
	})
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "skills.html", PageData{
		Title:    "Skills",
		Circles:  getCircles(),
		NavItems: getNavItems("/skills"),
	})
}

func experienceHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "experience.html", PageData{
		Title:    "Experience",
		Circles:  getCircles(),
		NavItems: getNavItems("/experience"),
	})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", PageData{
		Title:    "Contact",
		Circles:  getCircles(),
		NavItems: getNavItems("/contact"),
	})
}

func initPages() []Page {
	return []Page{
		{Route: "/", Title: "Home", Handler: indexHandler},
		{Route: "/projects", Title: "Projects", Handler: projectsHandler},
		{Route: "/skills", Title: "Skills", Handler: skillsHandler},
		{Route: "/experience", Title: "Experience", Handler: experienceHandler},
		{Route: "/contact", Title: "Contact", Handler: contactHandler},
	}
}

func getCircles() []BrokenCircle {
	return []BrokenCircle{
		{Diameter: "100", Rotate: -90},
		{Diameter: "90", Rotate: -120},
		{Diameter: "80", Rotate: -150},
		{Diameter: "70", Rotate: 180},
		{Diameter: "60", Rotate: 150},
		{Diameter: "50", Rotate: 120},
		{Diameter: "40", Rotate: 90},
		{Diameter: "30", Rotate: 60},
		{Diameter: "20", Rotate: 30},
		{Diameter: "10", Rotate: 0},
	}
}

func getNavItems(route string) []NavItem {
	navItems := []NavItem{}
	for i, p := range pages {
		navItems = append(navItems, NavItem{
			Title:  p.Title,
			Route:  p.Route,
			Active: route == p.Route,
			Last: i == len(pages)-1,
		})
	}
	return navItems
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
