package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BulkddataUrl          string         `json:"bulkdata_url"`
	DownloadDir           string         `json:"download_dir"`
	XMLDir                string         `json:"xml_dir"`
	JSONDir               string         `json:"json_dir"`
	JSONErrorDir          string         `json:"json_dir_error"`
	ExtractedDataDir      string         `json:"extracted_data_dir"`
	ExtractedFile         string         `json:"extracted_file"`
	CountryFilter         string         `json:"filter_for_country"`
	RunSteps              RunSteps       `json:"run_steps"`
	DeleteAfterProcessing bool           `json:"delete_after_processing"`
	DownloadPeriod        DownloadPeriod `json:"download_period"`
	LogLevel              string         `json:"log_level"`
}

type RunSteps struct {
	RunStepDownload   bool `json:"run_step_download"`
	RunStepProcessXML bool `json:"run_step_processxml"`
	RunStepTransform  bool `json:"run_step_transform"`
}

type DownloadPeriod struct {
	FromYear  int `json:"from_year"`
	FromMonth int `json:"from_month"`
	ToYear    int `json:"to_year"`
	ToMonth   int `json:"to_month"`
}

func LoadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
