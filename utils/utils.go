package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExtractWheel(filePath, destPath string) error {
	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		// Calculate the destination path and directory
		destPath := filepath.Join(destPath, file.Name)
		destDir := filepath.Dir(destPath)

		// Create the destination directory if it doesn't exist
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return err
		}

		// Open the file from the zip archive
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		// Create the destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy the contents from the zip file to the destination file
		_, err = io.Copy(destFile, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtractTarGz(src, dest string) error {
	fmt.Println(dest)
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return err
			}
		} else {
			file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_RDWR, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExtractTGZ(tgzFile string, destDir string) error {
	file, err := os.Open(tgzFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(destDir, header.Name)
		fileInfo := header.FileInfo()

		if fileInfo.IsDir() {
			if err := os.MkdirAll(targetFilePath, fileInfo.Mode()); err != nil {
				return err
			}
			continue
		}

		// Check if the file exists in the archive
		if !strings.HasPrefix(targetFilePath, destDir) {
			return fmt.Errorf("Attempted extraction outside of the destination directory: %s", targetFilePath)
		}

		outFile, err := os.Create(targetFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return err
		}
	}

	return nil
}

func RunGitleaks(directory string) {
	// Run the gitleaks command
	gitleaksCmd := exec.Command("gitleaks", "detect", "-v", "-s", directory, "--no-git", "-r", directory+"/output.json")
	gitleaksCmd.Stdout = os.Stdout
	gitleaksCmd.Stderr = os.Stderr

	if err := gitleaksCmd.Run(); err != nil {
		fmt.Println("Error running gitleaks:", err)
		os.Exit(1)
	}
}

// findTgzFiles returns a list of .tgz files in the specified directory.
func FindTgzFiles(directory string) ([]string, error) {
	var tgzFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tgz") {
			tgzFiles = append(tgzFiles, path)
		}
		return nil
	})

	return tgzFiles, err
}

// extractTgzFile extracts the contents of a .tgz file to the specified directory.
func ExtractTgzFile(tgzFile string, targetDirectory string) error {
	file, err := os.Open(tgzFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetDirectory, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories with appropriate permissions, including parent directories.
			if err := os.MkdirAll(targetFilePath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			// Create the target file's parent directory if it doesn't exist.
			parentDir := filepath.Dir(targetFilePath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return err
			}

			// Create the target file.
			file, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy the file contents.
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		default:
			// Skip unsupported file types.
			fmt.Printf("Skipped unsupported file type: %s\n", header.Typeflag)
		}
	}

	return nil
}
