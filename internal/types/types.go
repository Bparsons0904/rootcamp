package types

import "time"

type UserProgress struct {
	LessonID    string
	Completed   bool
	CompletedAt *time.Time
	Attempts    int
}

type Settings struct {
	SkipIntroAnimation bool
	UseBasicBash       bool
}

type SettingItem struct {
	Name        string
	DisplayName string
	Value       bool
	Description string
}

type FunFact struct {
	ID    string   `json:"id"`
	Tags  []string `json:"tags"`
	Title string   `json:"title"`
	Short string   `json:"short"`
	Full  string   `json:"full"`
}

type FunFactsData struct {
	Version string    `json:"version"`
	Facts   []FunFact `json:"facts"`
}

type Lesson struct {
	ID           string         `json:"id"`
	Command      string         `json:"command"`
	Code         string         `json:"code"`
	Title        string         `json:"title"`
	Tags         []string       `json:"tags"`
	Level        string         `json:"level"`
	Module       string         `json:"module"`
	About        LessonAbout    `json:"about"`
	Hints        []string       `json:"hints"`
	SkipSandbox  bool           `json:"skipSandbox,omitempty"`
	Sandbox      SandboxConfig  `json:"sandbox"`
	Instructions string         `json:"instructions"`
	Requirements []Requirement  `json:"requirements"`
}

type LessonAbout struct {
	What       string   `json:"what"`
	History    string   `json:"history"`
	Example    string   `json:"example"`
	CommonUses []string `json:"commonUses"`
}

type SandboxConfig struct {
	StartDir string            `json:"startDir"`
	Dirs     []string          `json:"dirs"`
	Files    map[string]string `json:"files"`
	Symlinks map[string]string `json:"symlinks"`
}

type Requirement struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Validator   string `json:"validator"`
	Expected    string `json:"expected"`
}

type LessonsData struct {
	Version string   `json:"version"`
	Lessons []Lesson `json:"lessons"`
}
