package loaders

import (
	"database/sql"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type PatchItem struct {
	Name string `json:"name"`
	Size int32  `json:"size"`
	Hash string `json:"hash"`
}

type PatchFile struct {
	Name      string      `json:"name"`
	PatchDate time.Time   `json:"patchDate"`
	Files     []PatchItem `json:"files"`
}

func (p *PatchFile) ToDomain() domain.Patch {
	var files []string = make([]string, len(p.Files))
	for idx, file := range p.Files {
		files[idx] = strings.ReplaceAll(file.Name, "\\", "/")
	}

	return domain.Patch{
		Name:   p.Name,
		Date:   p.PatchDate,
		Files:  files,
		Status: domain.PatchStatusPending,
	}
}

var sakrayPatchExpressions = [6]*regexp.Regexp{
	regexp.MustCompile(`(?i)^\d+-\d+-\d+_?rData`), // 9999-12-31rData.... or 9999-12-31_rData....
	regexp.MustCompile(`(?i)^\d+-\d+-\d+\wrData`), // 9999-12-31arData....
	regexp.MustCompile(`(?i)_sak[_\\.]`),          // 9999-12-31_foo_sak_a....
	regexp.MustCompile(`(?i)_sakra`),              // 9999-12-31_foo_sakra(y)....
	regexp.MustCompile(`(?i)DataSak`),             // 9999-12-31aDataSak....
	regexp.MustCompile(`(?i)ragexeRE`),            // 9999-12-31_ragExeRE.rgz
}

func (l *PatchListLoader) shouldSkipPatch(rawPatch PatchFile) bool {
	for _, exp := range sakrayPatchExpressions {
		if exp.Match([]byte(rawPatch.Name)) {
			return true
		}
	}

	if l.server.Type == server.ServerTypeKROMain {
		// 2020-02-25data_true_001.gpf is the first patch in kRO Server
		if rawPatch.PatchDate.After(time.Date(2020, 02, 24, 0, 0, 0, 0, time.UTC)) {
			return true
		}
	}

	return false
}

type PatchListLoader struct {
	repository *repository.PatchRepository
	server     *server.Server
}

func NewPatchListLoader(server *server.Server) *PatchListLoader {
	return &PatchListLoader{
		repository: server.Repositories.PatchRepository,
		server:     server,
	}
}

func (l *PatchListLoader) LoadFromJson(tx *sql.Tx, filePath string) error {
	var patchList []PatchFile
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &patchList)
	if err != nil {
		return err
	}

	// Note: old code ignored places where hash was empty; not sure if this matter

	for _, rawPatch := range patchList {
		if l.shouldSkipPatch(rawPatch) {
			continue
		}
		patch := rawPatch.ToDomain()

		err = l.repository.InsertPatch(tx, &patch)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Loads the initial set of patches from pre-processed JSON files.
 * These includes the patch name along with files list, in a order that _should_ be safe,
 * but not 100% garanteed to be the correct one. But it is the best we can get.
 *
 * This should only be used once, to populate the initial patch list.
 */
func (l *PatchListLoader) ExtractInitialPatchList(tx *sql.Tx) {
	files := [...]string{
		"../patches/kro/_plist/_init.json",
		"../patches/kro/_plist/2012.json",
		"../patches/kro/_plist/2013.json",
		"../patches/kro/_plist/2014.json",
		"../patches/kro/_plist/2015.json",
		"../patches/kro/_plist/2016.json",
		"../patches/kro/_plist/2017.json",
		"../patches/kro/_plist/2018.json",
		"../patches/kro/_plist/2019.json",
		"../patches/kro/_plist/2020.json",
		// "../patches/kro/_plist/2021.json",
		// "../patches/kro/_plist/2022.json",
		// "../patches/kro/_plist/2023.json",
		// "../patches/kro/_plist/2024.json",
		// "../patches/kro/_plist/2025.json",
	}

	for _, filePath := range files {
		if err := l.LoadFromJson(tx, filePath); err != nil {
			panic(err)
		}
	}
}
