package main

import (
	"fmt"
	"io/fs"
	"encoding/json"
)

/* == Types == */

type PageData struct {
	Title    string
	Projects []Project
	Circles  []BrokenCircle
	NavItems []NavItem
	EntityId int
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

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Repo        string   `json:"repo"`
	Details     string   `json:"details"`
}

/* == End Types == */


/* == Helper Funcs == */

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


/* == End Helper Funcs == */
