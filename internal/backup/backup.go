package backup

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Hash string
	Size int64
}

func Run(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			return err
		}

		fmt.Printf("%s %s\n", hash, path)

		compressFile(root)

		return nil
	})
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func compressFile(path string) error {
	srcFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %w", err)
	}
	defer srcFile.Close()

	tempDir, err := os.MkdirTemp("", "compression-")
	if err != nil {
		return fmt.Errorf("failed to create %w", err)
	}
	defer os.RemoveAll(tempDir)

	baseName := filepath.Base(path)
	zipPath := filepath.Join(tempDir, baseName + ".zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed create %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	archiveFile, err := zipWriter.Create(baseName)
	if err != nil {
		return fmt.Errorf("error archive %w", err)
	}

	_, err = io.Copy(archiveFile, srcFile)
	if err != nil {
		zipWriter.Close()
		return fmt.Errorf("error copy %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("zipwriter error %w", err)
	}

	fmt.Printf("succcess")
	return nil
}

func auth() {

}

func uploadFile(path string) {

}
