package main

import (
	"fmt"
	"os"
	"strings"
	"tedprocessor/config"
	"tedprocessor/convert"
	"tedprocessor/download"
	"tedprocessor/transform"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	zerolog.SetGlobalLevel(mapLogLevel(cfg.LogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting ETL process.")

	// Step 1: Download data
	// Ensure the destination directory exists
	if cfg.RunStepDownload {
		log.Debug().Msg("Starting Download process")

		err = os.MkdirAll(cfg.DownloadDir, 0755)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create destination directory")
			return
		}
		log.Debug().Str("directory", cfg.DownloadDir).Msg("Download directory exists")

		links, err := download.CreateDownloadLinks(cfg.BulkddataUrl, 0, 0, 0, 0)
		if err != nil {
			log.Error().Err(err).Msg("Error creating download links")
			return
		}

		for _, link := range links {
			err = download.DownloadAndExtract(link, cfg.DownloadDir)
			if err != nil {
				log.Error().Err(err).Msg("Error downloading data")
				return
			}
		}

		log.Info().Msg("Data downloaded successfully")
	}

	// Step 2: Convert to JSON
	// TODO itereate over all XML files
	if cfg.RunStepProcessXML {
		err = convert.ProcessXML("", cfg.DownloadDir, cfg.CountryFilter)
		if err != nil {
			log.Error().Err(err).Msg("Error reading xml data")
			return
		}
	}

	// Step 3: Build target data model and save
	if cfg.RunStepTransform {
		err = transform.ProcessData()
		if err != nil {
			log.Error().Err(err).Msg("Error processing json data to target data model")
			return
		}
	}

	log.Info().Msg("Processing finished. Exiting.")
}

// maps a string log level to a zerolog log level
func mapLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
