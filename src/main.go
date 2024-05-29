package main

import (
	"fmt"
	"os"
	"tedprocessor/config"
	"tedprocessor/convert"
	"tedprocessor/download"
	"tedprocessor/transform"
)

func main() {
	fmt.Println("Starting ETL process...")

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Step 1: Download data
	// Ensure the destination directory exists
	if cfg.RunStepDownload {
		err = os.MkdirAll(cfg.DownloadDir, 0755)
		if err != nil {
			fmt.Printf("Failed to create destination directory: %v", err)
			return
		}

		links, err := download.CreateDownloadLinks(cfg.BulkddataUrl, 0, 0, 0, 0)
		if err != nil {
			fmt.Printf("Error creating download links: %v\n", err)
			return
		}

		for _, link := range links {
			err = download.DownloadAndExtract(link, cfg.DownloadDir)
			if err != nil {
				fmt.Printf("Error downloading data: %v\n", err)
				return
			}
		}

		fmt.Println("Data downloaded successfully.")
	}

	// Step 2: Convert to JSON
	// TODO itereate over all XML files
	if cfg.RunStepProcessXML {
		err = convert.ProcessXML("", cfg.DownloadDir, cfg.CountryFilter)
		if err != nil {
			fmt.Printf("Error reading xml data: %v\n", err)
			return
		}
	}

	// Step 3: Build target data model and save
	if cfg.RunStepTransform {
		err = transform.ProcessData()
		if err != nil {
			fmt.Printf("Error processing json data to target model: %v\n", err)
			return
		}
	}

	fmt.Printf("Processing finished.")
}
