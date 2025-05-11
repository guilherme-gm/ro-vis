package domain

import (
	"errors"
	"strings"
	"time"
)

type UpdateChange struct {
	Patch string
	File  string
}

type Update struct {
	Date    time.Time
	Changes []UpdateChange
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

func (u Update) Name() string {
	return u.Date.Format(time.DateOnly)
}
