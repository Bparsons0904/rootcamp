package lab

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bobparsons/rootcamp/internal/types"
)

func ValidateRequirement(req types.Requirement, userInput, sandboxPath string) bool {
	switch req.Validator {
	case "exact":
		return validateExact(req.Expected, userInput)
	case "path_match":
		return validatePathMatch(req.Expected, userInput, sandboxPath)
	case "file_check":
		return validateFileExists(req.Expected, sandboxPath)
	case "regex":
		return validateRegex(req.Expected, userInput)
	default:
		return false
	}
}

func validateExact(expected, actual string) bool {
	return strings.TrimSpace(expected) == strings.TrimSpace(actual)
}

func validatePathMatch(expectedTemplate, actual, sandboxPath string) bool {
	expected := strings.ReplaceAll(expectedTemplate, "/tmp/rootcamp-{uuid}", sandboxPath)
	return strings.TrimSpace(expected) == strings.TrimSpace(actual)
}

func validateFileExists(filename, sandboxPath string) bool {
	fullPath := filepath.Join(sandboxPath, filename)
	_, err := os.Stat(fullPath)
	return err == nil
}

func validateRegex(pattern, actual string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(strings.TrimSpace(actual))
}

func ValidateLesson(lesson types.Lesson, userInput, sandboxPath string) (bool, string) {
	for _, req := range lesson.Requirements {
		if ValidateRequirement(req, userInput, sandboxPath) {
			return true, ""
		}
	}

	if len(lesson.Requirements) > 0 {
		return false, lesson.Requirements[0].Description
	}

	return false, "No requirements defined"
}
