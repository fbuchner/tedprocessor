package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BulkddataUrl        string `json:"bulkdata_url"`
	DataDirectory       string `json:"data_directory"`
	DownloadSubdir      string `json:"download_subdir"`
	XMLSubdir           string `json:"xml_subdir"`
	XMLErrorSubdir      string `json:"xml_error_subdir"`
	JSONSubdir          string `json:"json_subdir"`
	JSONErrorSubdir     string `json:"json_subdir_error"`
	ExtractedDataSubdir string `json:"extracted_data_subdir"`
	ExtractedFile       string `json:"extracted_file"`
	CountryFilter       string `json:"filter_for_country"`
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
