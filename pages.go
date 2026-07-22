package main

import (
	"fmt"
	"log"
	"net/http"
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

	entityId := 0
	var project Project
	var extraTemplates []string

	if idStr := r.URL.Query().Get("id"); idStr != "" {
		var id int
		if _, perr := fmt.Sscanf(idStr, "%d", &id); perr == nil {
			for _, p := range projects {
				if p.Id == id {
					entityId = id
					project = p
					extraTemplates = append(extraTemplates, p.Details)
					break
				}
			}
		}
	}

	renderTemplate(w, "projects.html", PageData{
		Title:          "Projects",
		Projects:       projects,
		Project:        project,
		Circles:        getCircles(),
		NavItems:       getNavItems("/projects"),
		EntityId:       entityId,
		ShowIcons8Link: true,
	}, extraTemplates...)
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	skills, err := loadSkills()
	if err != nil {
		log.Printf("load skills: %v", err)
	}
	featured, groups := groupSkills(skills)
	renderTemplate(w, "skills.html", PageData{
		Title:          "Skills",
		Circles:        getCircles(),
		NavItems:       getNavItems("/skills"),
		FeaturedSkills: featured,
		SkillGroups:    groups,
		ShowIcons8Link: true,
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
