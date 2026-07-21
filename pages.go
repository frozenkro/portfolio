package main

import (
	"net/http"
	"log"
)

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
