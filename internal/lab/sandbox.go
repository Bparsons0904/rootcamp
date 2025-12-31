package lab

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bobparsons/rootcamp/internal/types"
)

func generateShortID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 5

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[rng.Intn(len(charset))]
	}
	return string(id)
}

func Create(lesson types.Lesson) (string, error) {
	sandboxID := generateShortID()
	sandboxPath := filepath.Join("/tmp", fmt.Sprintf("rootcamp-%s", sandboxID))

	if err := os.MkdirAll(sandboxPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create sandbox: %w", err)
	}

	for _, dir := range lesson.Sandbox.Dirs {
		dirPath := filepath.Join(sandboxPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	for filePath, content := range lesson.Sandbox.Files {
		fullPath := filepath.Join(sandboxPath, filePath)

		dirPath := filepath.Dir(fullPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
	}

	for linkName, target := range lesson.Sandbox.Symlinks {
		linkPath := filepath.Join(sandboxPath, linkName)

		dirPath := filepath.Dir(linkPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create parent directory for symlink: %w", err)
		}

		if err := os.Symlink(target, linkPath); err != nil {
			Cleanup(sandboxPath)
			return "", fmt.Errorf("failed to create symlink %s: %w", linkName, err)
		}
	}

	return sandboxPath, nil
}

func GetStartPath(sandboxPath string, lesson types.Lesson) string {
	if lesson.Sandbox.StartDir == "" {
		return sandboxPath
	}
	return filepath.Join(sandboxPath, lesson.Sandbox.StartDir)
}

func Cleanup(sandboxPath string) error {
	if sandboxPath == "" || !strings.HasPrefix(sandboxPath, "/tmp/rootcamp-") {
		return fmt.Errorf("invalid sandbox path: %s", sandboxPath)
	}
	return os.RemoveAll(sandboxPath)
}
