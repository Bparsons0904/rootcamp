package lessons

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bobparsons/rootcamp/internal/types"
)

//go:embed data/lessons/*.json
var embeddedLessonsFS embed.FS

var cachedLessons *types.LessonsData

func LoadLessons() (*types.LessonsData, error) {
	if cachedLessons != nil {
		return cachedLessons, nil
	}

	data := &types.LessonsData{
		Version: "1.0",
		Lessons: []types.Lesson{},
	}

	entries, err := embeddedLessonsFS.ReadDir("data/lessons")
	if err != nil {
		return nil, fmt.Errorf("failed to read lessons directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			if err := loadLessonsFromFile(entry.Name(), data); err != nil {
				return nil, fmt.Errorf("failed to load %s: %w", entry.Name(), err)
			}
		}
	}

	cachedLessons = data
	return cachedLessons, nil
}

func loadLessonsFromFile(filename string, accumulator *types.LessonsData) error {
	filePath := filepath.Join("data/lessons", filename)
	content, err := embeddedLessonsFS.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	var fileData types.LessonsData
	if err := json.Unmarshal(content, &fileData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	accumulator.Lessons = append(accumulator.Lessons, fileData.Lessons...)
	return nil
}

func GetLessonByID(id string) (*types.Lesson, error) {
	data, err := LoadLessons()
	if err != nil {
		return nil, err
	}

	for i := range data.Lessons {
		if data.Lessons[i].ID == id {
			return &data.Lessons[i], nil
		}
	}

	return nil, fmt.Errorf("lesson not found: %s", id)
}

func GetLessonsByTag(tag string) ([]types.Lesson, error) {
	data, err := LoadLessons()
	if err != nil {
		return nil, err
	}

	var filtered []types.Lesson
	for _, lesson := range data.Lessons {
		if slices.Contains(lesson.Tags, tag) {
			filtered = append(filtered, lesson)
		}
	}

	return filtered, nil
}

func GetLessonsByLevel(level string) ([]types.Lesson, error) {
	data, err := LoadLessons()
	if err != nil {
		return nil, err
	}

	var filtered []types.Lesson
	for _, lesson := range data.Lessons {
		if lesson.Level == level {
			filtered = append(filtered, lesson)
		}
	}

	return filtered, nil
}

func GetLessonsByModule(module string) ([]types.Lesson, error) {
	data, err := LoadLessons()
	if err != nil {
		return nil, err
	}

	var filtered []types.Lesson
	for _, lesson := range data.Lessons {
		if lesson.Module == module {
			filtered = append(filtered, lesson)
		}
	}

	return filtered, nil
}

func GetAllModules() ([]string, error) {
	data, err := LoadLessons()
	if err != nil {
		return nil, err
	}

	moduleSet := make(map[string]bool)
	for _, lesson := range data.Lessons {
		if lesson.Module != "" {
			moduleSet[lesson.Module] = true
		}
	}

	modules := make([]string, 0, len(moduleSet))
	for module := range moduleSet {
		modules = append(modules, module)
	}

	return modules, nil
}
