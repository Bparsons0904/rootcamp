package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bobparsons/rootcamp/internal/types"
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

	if err := InitDefaultSettings(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize default settings: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS settings (
		setting_name TEXT PRIMARY KEY,
		setting_value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

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

func InitDefaultSettings(db *sql.DB) error {
	defaults := map[string]string{
		"skip_intro_animation": "false",
	}

	for name, value := range defaults {
		query := `
			INSERT OR IGNORE INTO settings (setting_name, setting_value)
			VALUES (?, ?)
		`
		_, err := db.Exec(query, name, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSetting(db *sql.DB, name string, defaultValue string) (string, error) {
	query := `SELECT setting_value FROM settings WHERE setting_name = ?`

	var value string
	err := db.QueryRow(query, name).Scan(&value)

	if err == sql.ErrNoRows {
		return defaultValue, nil
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func SetSetting(db *sql.DB, name string, value string) error {
	query := `
		INSERT INTO settings (setting_name, setting_value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(setting_name) DO UPDATE SET
			setting_value = ?,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.Exec(query, name, value, value)
	return err
}

func GetAllSettings(db *sql.DB) (*types.Settings, error) {
	skipIntro, err := GetSetting(db, "skip_intro_animation", "false")
	if err != nil {
		return nil, err
	}

	return &types.Settings{
		SkipIntroAnimation: stringToBool(skipIntro),
	}, nil
}

func stringToBool(s string) bool {
	return s == "true"
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
