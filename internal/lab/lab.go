package lab

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/Bparsons0904/rootcamp/internal/types"
)

func Create(lesson types.Lesson) (string, error) {
	sandboxID := uuid.New().String()
	sandboxPath := filepath.Join("/tmp", fmt.Sprintf("rootcamp-%s", sandboxID))

	if err := os.MkdirAll(sandboxPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create sandbox directory: %w", err)
	}

	for _, dir := range lesson.Lab.Dirs {
		dirPath := filepath.Join(sandboxPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	for filePath, content := range lesson.Lab.Files {
		fullPath := filepath.Join(sandboxPath, filePath)

		dirPath := filepath.Dir(fullPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create parent directory for %s: %w", filePath, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
	}

	return sandboxPath, nil
}

func Cleanup(sandboxPath string) error {
	if sandboxPath == "" || !strings.HasPrefix(sandboxPath, "/tmp/rootcamp-") {
		return fmt.Errorf("invalid sandbox path: %s", sandboxPath)
	}

	return os.RemoveAll(sandboxPath)
}

func Validate(userInput, secretCode string) bool {
	return strings.TrimSpace(userInput) == strings.TrimSpace(secretCode)
}
