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
