package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
)

/* == Types == */

type PageData struct {
	Title          string
	Projects       []Project
	Circles        []BrokenCircle
	NavItems       []NavItem
	EntityId       int
	ShowIcons8Link bool
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

type Tag struct {
	Name string
	Src  string
	Alt  string
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	t.Name = name
	if icon, ok := tagIcons[name]; ok {
		t.Src = icon.Src
		t.Alt = icon.Alt
	}
	return nil
}

type Project struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags"`
	Repo        string `json:"repo"`
	Details     string `json:"details"`
}

/* == End Types == */

/* == Conststs == */

var tagIcons = map[string]struct{ Src, Alt string }{
	"Go":         {"/images/icons8-golang.svg", "Go logo"},
	"k3s":        {"/images/icons8-kubernetes.svg", "k3s logo"},
	"MQTT":       {"/images/Mqtt-hor.svg", "MQTT logo"},
	"Grafana":    {"/images/icons8-grafana.svg", "Grafana logo"},
	"C":          {"/images/icons8-c.svg", "C logo"},
	"I2C":        {"/images/i2c_bus_logo.svg", "I2C logo"},
	"MCP":        {"/images/Model_Context_Protocol_logo.svg", "MCP logo"},
	"Zig":        {"/images/Zig_logo_2020.svg", "Zig logo"},
	"wayland":    {"/images/Wayland_Logo.svg", "Wayland logo"},
	"PostgreSQL": {"/images/icons8-postgresql-48.png", "PostgreSQL logo"},
	"Docker":     {"/images/icons8-docker-48.png", "Docker logo"},
	"CMake":      {"/images/icons8-cmake-48.png", "CMake logo"},
	"Ollama":     {"/images/ollama.png", "Ollama logo"},
	"Ubuntu":     {"/images/icons8-ubuntu-48.png", "Ubuntu logo"},
	"Kotlin":     {"/images/icons8-kotlin-48.png", "Kotlin logo"},
	"Android":    {"/images/icons8-android-48.png", "Android logo"},
	"Ansible":    {"/images/icons8-ansible-48.png", "Ansible logo"},
	"Bash":       {"/images/icons8-bash.svg", "Bash logo"},
	"Claude":     {"/images/icons8-claude-48.png", "Claude logo"},
}

/* == End Csnosonsnsts == */

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
			Last:   i == len(pages)-1,
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
