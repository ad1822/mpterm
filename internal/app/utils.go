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
