package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer fromFile.Close()

	fileInfo, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek in source file: %w", err)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer toFile.Close()

	if limit == 0 || limit > fileInfo.Size()-offset {
		limit = fileInfo.Size() - offset
	}

	buf := make([]byte, 1024)
	var copied int64

	for copied < limit {
		bytesToRead := int64(len(buf))
		if limit-copied < bytesToRead {
			bytesToRead = limit - copied
		}

		n, err := fromFile.Read(buf[:bytesToRead])
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("could not read from source file: %w", err)
		}

		if n == 0 {
			break
		}

		_, err = toFile.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("could not write to destination file: %w", err)
		}

		copied += int64(n)

		// Print progress
		fmt.Printf("\rCopying... %d%% complete", (copied*100)/limit)
	}

	fmt.Println("\nCopy complete")
	return nil
}
