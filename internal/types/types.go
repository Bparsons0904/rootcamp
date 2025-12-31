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
