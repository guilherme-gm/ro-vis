package dao

import (
	"encoding/json"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func (p *Patch) ToDomain() domain.Patch {
	var files []string
	err := json.Unmarshal(p.Files, &files)
	if err != nil {
		panic(err)
	}

	return domain.Patch{
		Id:     p.ID,
		Name:   p.Name,
		Date:   p.Date,
		Files:  files,
		Status: domain.PatchStatus(p.Status),
	}
}
