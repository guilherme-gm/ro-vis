package domain

import (
	"fmt"
	"regexp"
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

	return UpdateChange{}, NewNotFoundError(fmt.Sprintf("file '%s' not found in update", file))
}

func (u Update) GetChangesForFile(file *regexp.Regexp) ([]UpdateChange, error) {
	var changes []UpdateChange
	for _, change := range u.Changes {
		if file.MatchString(change.File) {
			changes = append(changes, change)
		}
	}

	if len(changes) == 0 {
		return changes, NewNotFoundError(fmt.Sprintf("file '%s' not found in update", file))
	}

	return changes, nil
}

func (u Update) HasChangedAnyFiles(files []*regexp.Regexp) bool {
	for _, change := range u.Changes {
		for _, targetFile := range files {
			if targetFile.MatchString(change.File) {
				return true
			}
		}
	}

	return false
}

func (u Update) Name() string {
	return u.Date.Format(time.DateOnly)
}
