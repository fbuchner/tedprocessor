package main

import (
	"fmt"
	"os"
	"tedprocessor/config"
	"tedprocessor/download"
)

func main() {
	fmt.Println("Starting ETL process...")

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(cfg.DownloadSubdir, 0755); err != nil {
		fmt.Printf("Failed to create destination directory: %v", err)
	}

	// Step 1: Download data
	err = download.DownloadAndExtract(cfg.BulkddataUrl, cfg.DownloadSubdir)
	if err != nil {
		fmt.Printf("Error downloading data: %v\n", err)
		return
	}
	fmt.Println("Data downloaded successfully.")
}
