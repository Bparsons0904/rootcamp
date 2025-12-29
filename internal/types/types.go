package types

import "time"

type Lesson struct {
	ID          string
	Title       string
	Order       int
	What        string
	Example     string
	History     string
	CommonUses  []string
	Lab         LabConfig
	SecretCode  string
	Difficulty  string
	Tags        []string
}

type LabConfig struct {
	Dirs  []string
	Files map[string]string
	Task  string
}

type UserProgress struct {
	LessonID    string
	Completed   bool
	CompletedAt *time.Time
	Attempts    int
}
