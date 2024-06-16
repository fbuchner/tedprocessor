package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tedprocessor/config"
	"tedprocessor/convert"
	"tedprocessor/download"
	"tedprocessor/transform"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// The path to the configuration file can be given as a parameter ("go run main.go -config=/path/to/your/config.json")
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path to the config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	//Set log level according to configuration, defaults to info (see function below)
	zerolog.SetGlobalLevel(mapLogLevel(cfg.LogLevel))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting ETL process.")

	// Ensure the destination directory exists
	dirs := []string{cfg.DownloadDir, cfg.XMLDir, cfg.JSONDir}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to create destination directory: %s", dir)
			return
		}
		log.Debug().Str("directory", dir).Msg("Found or created directory")
	}

	// Step 1: Download data
	if cfg.RunSteps.RunStepDownload {
		log.Info().Msg("Starting Download process.")

		links, err := download.CreateDownloadLinks(cfg.BulkddataUrl, cfg.DownloadPeriod.FromYear, cfg.DownloadPeriod.FromMonth, cfg.DownloadPeriod.ToYear, cfg.DownloadPeriod.ToMonth)
		if err != nil {
			log.Error().Err(err).Msg("Error creating download links")
			return
		}
		log.Debug().Interface("Download links", links).Msg("Found download links")

		for _, link := range links {
			log.Debug().Str("URL", link).Msg("Downloading from link")
			err = download.DownloadAndExtract(link, cfg.DownloadDir, cfg.XMLDir)
			if err != nil {
				log.Error().Err(err).Msg("Error downloading data")
				return
			}
		}

		log.Info().Msg("Data downloaded successfully.")
	}

	// Step 2: Convert to JSON
	if cfg.RunSteps.RunStepProcessXML {
		log.Info().Msg("Starting processing of XML files.")

		//iterate directory tree
		err = filepath.Walk(cfg.XMLDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
				err = convert.ProcessXML(path, cfg.JSONDir, cfg.CountryFilter, cfg.DeleteAfterProcessing)
				if err != nil {
					return err
				}
				if cfg.DeleteAfterProcessing {
					err := convert.RemoveFile(path)
					if err != nil {
						return fmt.Errorf("failed to remove xml file: %v", err)
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Error().Err(err).Msg("Error reading xml data")
			return
		}

		log.Info().Msg("XML files processed successfully.")
	}

	// Step 3: Build target data model and save
	if cfg.RunSteps.RunStepTransform {
		log.Info().Msg("Starting transformation to target model for output.")

		err = filepath.Walk(cfg.JSONDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
				err = transform.ProcessJSON(path, cfg.ExtractedFile, cfg.CSVSeparator, cfg.DeleteAfterProcessing)
				if err != nil {
					return err
				}
				if cfg.DeleteAfterProcessing {
					err := transform.RemoveFile(path)
					if err != nil {
						return fmt.Errorf("failed to remove json file: %v", err)
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Error().Err(err).Msg("Error processing JSON data")
			return
		}

		log.Info().Msg("Data transformation finished successfully.")
	}

	log.Info().Msg("Processing finished. Exiting.")
}

// Maps a string log level to a zerolog log level
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
