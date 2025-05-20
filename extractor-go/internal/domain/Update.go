package domain

import (
	"errors"
	"strings"
	"time"
)

type UpdateChange struct {
	Patch string `json:"patch"`
	File  string `json:"file"`
}

type Update struct {
	Date    time.Time      `json:"date"`
	Changes []UpdateChange `json:"changes"`
}

func (u Update) GetChangeForFile(file string) (UpdateChange, error) {
	file = strings.ToLower(file)
	for _, change := range u.Changes {
		if strings.ToLower(change.File) == file {
			return change, nil
		}
	}

	return UpdateChange{}, errors.New("could not find file - " + file)
}

func (u Update) HasChangedAnyFiles(files []string) bool {
	for idx, file := range files {
		files[idx] = strings.ToLower(file)
	}

	for _, change := range u.Changes {
		for _, targetFile := range files {
			if strings.ToLower(change.File) == targetFile {
				return true
			}
		}
	}

	return false
}

func (u Update) Name() string {
	return u.Date.Format(time.DateOnly)
}
