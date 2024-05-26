package download

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CreateDownloadLinks(downloadUrl string, startYear, startMonth, endYear, endMonth int) ([]string, error) {
	// Validate months
	if startMonth < 1 || startMonth > 12 {
		return nil, errors.New("start month must be between 1 and 12")
	}
	if endMonth < 1 || endMonth > 12 {
		return nil, errors.New("end month must be between 1 and 12")
	}

	// Create time instances for the start and end dates
	startDate := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	var urls []string

	// Iterate over each month between the start date and the end date
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 1, 0) {
		// Format the URL
		formattedDate := fmt.Sprintf("%d-%02d", date.Year(), date.Month())
		url := fmt.Sprintf("%s/%s", downloadUrl, formattedDate)
		// Append the URL to the slice
		urls = append(urls, url)
	}

	return urls, nil
}

// ExtractTarGz extracts a .tar.gz file to the specified directory.
func ExtractTarGz(gzipPath, destDir string) error {
	file, err := os.Open(gzipPath)
	if err != nil {
		return fmt.Errorf("failed to open gzip file: %v", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %v", err)
		}

		filePath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		case tar.TypeReg:
			// Ensure the directory exists
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}

			outFile, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to copy file contents: %v", err)
			}
			outFile.Close()

			if err := os.Chmod(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set file permissions: %v", err)
			}
		default:
			return fmt.Errorf("unsupported file type: %v", header.Typeflag)
		}
	}

	return nil
}

// DownloadAndExtract downloads a .tar.gz file from the given URL, saves it to disk, and extracts it to the specified directory.
func DownloadAndExtract(url, destDir string) error {
	tmpFile, err := os.CreateTemp("", "download-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := DownloadFile(url, tmpFile.Name()); err != nil {
		return err
	}

	if err := ExtractTarGz(tmpFile.Name(), destDir); err != nil {
		return err
	}

	return nil
}
