package lessons

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/bobparsons/rootcamp/internal/types"
)

//go:embed data/course.json
var embeddedCourseFS embed.FS

var cachedCourse *types.CourseData

func LoadCourse() (*types.CourseData, error) {
	if cachedCourse != nil {
		return cachedCourse, nil
	}

	content, err := embeddedCourseFS.ReadFile("data/course.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read course.json: %w", err)
	}

	var courseData types.CourseData
	if err := json.Unmarshal(content, &courseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal course: %w", err)
	}

	cachedCourse = &courseData
	return cachedCourse, nil
}

func GetCourseLessons(progressMap map[string]*types.UserProgress, allLessons []types.Lesson) ([]types.CourseLessonItem, error) {
	courseData, err := LoadCourse()
	if err != nil {
		return nil, err
	}

	lessonMap := make(map[string]*types.Lesson)
	for i := range allLessons {
		lessonMap[allLessons[i].ID] = &allLessons[i]
	}

	var courseLessons []types.CourseLessonItem

	for idx, courseLessonRef := range courseData.Course.Lessons {
		lesson, exists := lessonMap[courseLessonRef.LessonId]
		if !exists {
			continue
		}

		status := types.LessonLocked

		if idx == 0 {
			status = types.LessonUnlocked
		} else if idx > 0 {
			prevLesson := &courseLessons[idx-1]
			if prevLesson.Status == types.LessonComplete {
				status = types.LessonUnlocked
			}
		}

		if progress, exists := progressMap[lesson.ID]; exists && progress.Completed {
			status = types.LessonComplete
		}

		courseLessons = append(courseLessons, types.CourseLessonItem{
			Lesson:   *lesson,
			Status:   status,
			Sequence: courseLessonRef.Sequence,
		})
	}

	return courseLessons, nil
}

func GetNextUnlockedLesson(courseLessons []types.CourseLessonItem) *types.CourseLessonItem {
	for i := range courseLessons {
		if courseLessons[i].Status == types.LessonUnlocked {
			return &courseLessons[i]
		}
	}
	return nil
}

func GetCourseProgress(courseLessons []types.CourseLessonItem) (int, int) {
	completed := 0
	for _, item := range courseLessons {
		if item.Status == types.LessonComplete {
			completed++
		}
	}
	return completed, len(courseLessons)
}
