package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFullPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s\n", err)
	}
	fullPath := filepath.Join(homeDir + "/Music")
	return fullPath
}

func GetQueueDBPath() string {
	dir, _ := os.UserHomeDir()
	dir = filepath.Join(dir, ".local", "share")
	dbDir := filepath.Join(dir, "mpterm")
	os.MkdirAll(dbDir, 0755)
	return filepath.Join(dbDir, "queue.db")
}
