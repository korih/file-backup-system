package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Run(ctx context.Context, folder, bucket, gpgRecipient string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	tmpFile, err := os.CreateTemp("", "homelab-backup-*.tar.gz.gpg")
	if err != nil {
		return fmt.Errorf("create temporary backup file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if err := createdEncryptedArchive(folder, tmpFile); err != nil {
		return fmt.Errorf("create backup archive: %w", err)
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind backup file: %w", err)
	}

	datePath := time.Now().Format("2003/01/02")
	
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(datePath + "/backup.tar.gz.gpg"),
		Body:   tmpFile,
	})

	if err != nil {
		return fmt.Errorf("upload to S3: %w", err)
	}

	return nil
}

func createdEncryptedArchive(folder string, output io.Writer) error {
	cmd := exec.Command(
		"gpg",
		"--batch",
		"--yes",
		"--encrypt",
		"--recipient",
		"kori@korih.com",
		"--output",
		"-",
	);

	gpgInput, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error with pipe: %w", err)
	}

	cmd.Stdout = output
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error with gpg %w", err)
	}

	gzipWriter := gzip.NewWriter(gpgInput)

	tarWriter := tar.NewWriter(gzipWriter)

	err = addFolderToTar(tarWriter, folder)

	if closeErr := tarWriter.Close(); err == nil {
		err = closeErr
	}

	if closeErr := gzipWriter.Close(); err == nil {
		err = closeErr
	}

	if closeErr := gpgInput.Close(); err == nil {
		err = closeErr
	}

	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return fmt.Errorf("close error %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait error %w", err)
	}

	return nil
}

func addFolderToTar(tw *tar.Writer, folder string) error {
	return filepath.Walk(folder, func (path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error %w", err)
		}

		relativePath, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = relativePath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
}
