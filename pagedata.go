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
	Project        Project
	Circles        []BrokenCircle
	NavItems       []NavItem
	EntityId       int
	ShowIcons8Link bool
	FeaturedSkills []Skill
	SkillGroups    []SkillGroup
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
	Name     string
	Src      string
	Alt      string
	NameIcon bool // Some icons have the name in them
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
		t.NameIcon = icon.NameIcon
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

type Skill struct {
	Name        string   `json:"name"`
	Cat         string   `json:"cat"`
	Featured    bool     `json:"featured"`
	Description string   `json:"description"`
	Src         string   `json:"-"`
	Alt         string   `json:"-"`
	NameIcon    bool     `json:"-"`
	Category    SkillCat `json:"-"`
}

type SkillGroup struct {
	Category SkillCat
	Skills   []Skill
}

type SkillCat string

const (
	BackendWeb      SkillCat = "Backend Web"
	Devops          SkillCat = "Devops"
	AgenticAI       SkillCat = "Tooling and AI"
	EmbeddedSystems SkillCat = "Embedded / Systems"
	FrontendWeb     SkillCat = "Frontend Web"
	WishIKnewLess   SkillCat = "Things I Wish I Knew Less About"
)

/* == End Types == */

/* == Conststs == */

var tagIcons = map[string]struct {
	Src, Alt string
	NameIcon bool
}{
	"Go":         {"/images/icons8-golang.svg", "Go logo", false},
	"k3s":        {"/images/icons8-kubernetes.svg", "k3s logo", false},
	"MQTT":       {"/images/Mqtt-hor.svg", "MQTT logo", true},
	"Grafana":    {"/images/icons8-grafana.svg", "Grafana logo", false},
	"C":          {"/images/icons8-c.svg", "C logo", false},
	"I2C":        {"/images/i2c_bus_logo.svg", "I2C logo", false},
	"MCP":        {"/images/Model_Context_Protocol_logo.svg", "MCP logo", false},
	"Zig":        {"/images/Zig_logo_2020.svg", "Zig logo", true},
	"wayland":    {"/images/Wayland_Logo.svg", "Wayland logo", false},
	"PostgreSQL": {"/images/icons8-postgresql-48.png", "PostgreSQL logo", false},
	"Docker":     {"/images/icons8-docker-48.png", "Docker logo", false},
	"CMake":      {"/images/icons8-cmake-48.png", "CMake logo", false},
	"Ollama":     {"/images/ollama.png", "Ollama logo", false},
	"Ubuntu":     {"/images/icons8-ubuntu-48.png", "Ubuntu logo", false},
	"Kotlin":     {"/images/icons8-kotlin-48.png", "Kotlin logo", false},
	"Android":    {"/images/icons8-android-48.png", "Android logo", false},
	"Ansible":    {"/images/icons8-ansible-48.png", "Ansible logo", false},
	"Bash":       {"/images/icons8-bash.svg", "Bash logo", false},
	"Claude":     {"/images/icons8-claude-48.png", "Claude logo", false},
}

var skillCats = []SkillCat{BackendWeb, Devops, AgenticAI, EmbeddedSystems, FrontendWeb, WishIKnewLess}

// catKeys maps JSON "cat" values to their SkillCat display name.
var catKeys = map[string]SkillCat{
	"backendWeb":    BackendWeb,
	"devops":        Devops,
	"agenticAI":     AgenticAI,
	"embedded":      EmbeddedSystems,
	"frontendWeb":   FrontendWeb,
	"wishIKnewLess": WishIKnewLess,
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

func loadSkills() ([]Skill, error) {
	data, err := fs.ReadFile(contentFS, "content/skills.json")
	if err != nil {
		return nil, fmt.Errorf("read skills.json: %w", err)
	}
	var skills []Skill
	if err := json.Unmarshal(data, &skills); err != nil {
		return nil, fmt.Errorf("unmarshal skills.json: %w", err)
	}
	for i := range skills {
		if icon, ok := tagIcons[skills[i].Name]; ok {
			skills[i].Src = icon.Src
			skills[i].Alt = icon.Alt
			skills[i].NameIcon = icon.NameIcon
		}
		if cat, ok := catKeys[skills[i].Cat]; ok {
			skills[i].Category = cat
		}
	}
	return skills, nil
}

func groupSkills(skills []Skill) ([]Skill, []SkillGroup) {
	var featured []Skill
	groups := make(map[SkillCat][]Skill)

	for _, s := range skills {
		if s.Featured {
			featured = append(featured, s)
		}
		if s.Category == "" {
			continue
		}
		groups[s.Category] = append(groups[s.Category], s)
	}

	var result []SkillGroup
	for _, cat := range skillCats {
		if skills, ok := groups[cat]; ok {
			result = append(result, SkillGroup{Category: cat, Skills: skills})
		}
	}
	return featured, result
}

/* == End Helper Funcs == */
