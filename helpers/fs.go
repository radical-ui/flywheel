package helpers

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

func EnsureDirExists(dir string) error {
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return nil
	}

	if err == nil {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.Join(fmt.Errorf("failed to create directory '%s' after checking to make sure it doesn't exist", dir), err)
	}

	return nil
}

func GetLatestModifiedTime(dir string) (time.Time, error) {
	var latestTime time.Time

	files, err := os.ReadDir(dir)
	if err != nil {
		return latestTime, err
	}

	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			return latestTime, err
		}

		entryModTime := info.ModTime()
		if entryModTime.After(latestTime) {
			latestTime = entryModTime
		}

		if info.IsDir() {
			dirModTime, err := GetLatestModifiedTime(path.Join(dir, entry.Name()))
			if err != nil {
				return latestTime, err
			}

			if dirModTime.After(latestTime) {
				latestTime = dirModTime
			}
		}
	}

	return latestTime, nil
}
