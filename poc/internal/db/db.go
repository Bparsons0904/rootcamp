package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Bparsons0904/rootcamp/internal/types"
)

func InitDB() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".rootcamp")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .rootcamp directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "rootcamp.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS progress (
		lesson_id TEXT PRIMARY KEY,
		completed BOOLEAN DEFAULT FALSE,
		completed_at DATETIME,
		attempts INTEGER DEFAULT 0
	);
	`

	_, err := db.Exec(schema)
	return err
}

func GetProgress(db *sql.DB, lessonID string) (*types.UserProgress, error) {
	query := `
		SELECT lesson_id, completed, completed_at, attempts
		FROM progress
		WHERE lesson_id = ?
	`

	var progress types.UserProgress
	var completedAt sql.NullTime

	err := db.QueryRow(query, lessonID).Scan(
		&progress.LessonID,
		&progress.Completed,
		&completedAt,
		&progress.Attempts,
	)

	if err == sql.ErrNoRows {
		return &types.UserProgress{
			LessonID:    lessonID,
			Completed:   false,
			CompletedAt: nil,
			Attempts:    0,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		progress.CompletedAt = &completedAt.Time
	}

	return &progress, nil
}

func GetAllProgress(db *sql.DB) (map[string]*types.UserProgress, error) {
	query := `SELECT lesson_id, completed, completed_at, attempts FROM progress`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progressMap := make(map[string]*types.UserProgress)

	for rows.Next() {
		var progress types.UserProgress
		var completedAt sql.NullTime

		err := rows.Scan(
			&progress.LessonID,
			&progress.Completed,
			&completedAt,
			&progress.Attempts,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			progress.CompletedAt = &completedAt.Time
		}

		progressMap[progress.LessonID] = &progress
	}

	return progressMap, rows.Err()
}

func MarkComplete(db *sql.DB, lessonID string) error {
	query := `
		INSERT INTO progress (lesson_id, completed, completed_at, attempts)
		VALUES (?, TRUE, ?, 0)
		ON CONFLICT(lesson_id) DO UPDATE SET
			completed = TRUE,
			completed_at = ?
	`

	now := time.Now()
	_, err := db.Exec(query, lessonID, now, now)
	return err
}

func IncrementAttempts(db *sql.DB, lessonID string) error {
	query := `
		INSERT INTO progress (lesson_id, completed, completed_at, attempts)
		VALUES (?, FALSE, NULL, 1)
		ON CONFLICT(lesson_id) DO UPDATE SET
			attempts = attempts + 1
	`

	_, err := db.Exec(query, lessonID)
	return err
}
