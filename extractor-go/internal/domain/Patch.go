package domain

import (
	"fmt"
	"regexp"
	"time"
)

const (
	PatchStatusPending   = "pending"
	PatchStatusExtracted = "extracted"
	PatchStatusGone      = "gone"
	PatchStatusSkipped   = "skipped"
)

type PatchStatus string

func (ps PatchStatus) String() string {
	return string(ps)
}

// A patch represents a single update container, containing new files.
// Sometimes, we get "Updates", which are a set of patch files applied in the same day,
// e.g. due to a weekly maintenance.
type Patch struct {
	Id     int32
	Name   string
	Date   time.Time
	Files  []string
	Status PatchStatus
}

var patchDatePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`), // 2024-10-30blablabla
	regexp.MustCompile(`(\d{4})(\d{2})(\d{2})`),       // 20241030blablabla
	regexp.MustCompile(`(\d{4})-(\d{2})(\d{2})`),      // 2024-1030blablabla
}

func TryGetPatchDate(patchName string) time.Time {
	for _, pattern := range patchDatePatterns {
		dateParts := pattern.FindStringSubmatch(patchName)
		if dateParts == nil || len(dateParts) != 4 {
			continue
		}

		dateStr := fmt.Sprintf("%s-%s-%s", dateParts[1], dateParts[2], dateParts[3])
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		return date
	}

	return time.Time{}
}
