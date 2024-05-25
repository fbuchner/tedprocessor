package download

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file from the given URL and saves it to the specified file path.
func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %v", err)
	}

	return nil
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
