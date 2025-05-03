package domain

import "time"

type Patch struct {
	Id    int32
	Name  string
	Date  time.Time
	Files []string
}
