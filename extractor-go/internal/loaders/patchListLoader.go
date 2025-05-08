package loaders

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
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
		Name:  p.Name,
		Date:  p.PatchDate,
		Files: files,
	}
}

var sakrayPatchExpressions = [4]*regexp.Regexp{
	regexp.MustCompile(`(?i)^\d+-\d+-\d+rData`), // 9999-12-31rData....
	regexp.MustCompile(`(?i)_sak[_\.]`),         // 9999-12-31_foo_sak_a....
	regexp.MustCompile(`(?i)_sakra`),            // 9999-12-31_foo_sakra(y)....
	regexp.MustCompile(`(?i)ragexeRE`),          // 9999-12-31_ragExeRE.rgz
}

func shouldSkipPatch(rawPatch PatchFile) bool {
	for _, exp := range sakrayPatchExpressions {
		if exp.Match([]byte(rawPatch.Name)) {
			return true
		}
	}

	// TODO: Remove this -- temporary logic for testing
	if rawPatch.PatchDate.Before(time.Date(2018, time.March, 21, 0, 0, 0, 0, time.UTC)) {
		return true
	}

	return false
}

func loadFromJson(filePath string) error {
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

	patchRepository := repository.GetPatchRepository()
	for _, rawPatch := range patchList {
		if shouldSkipPatch(rawPatch) {
			continue
		}
		patch := rawPatch.ToDomain()

		err = patchRepository.InsertPatch(&patch)
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
func ExtractInitialPatchList() {
	files := [...]string{
		"../patches/_plist/_init.json",
		"../patches/_plist/2012.json",
		"../patches/_plist/2013.json",
		"../patches/_plist/2014.json",
		"../patches/_plist/2015.json",
		"../patches/_plist/2016.json",
		"../patches/_plist/2017.json",
		"../patches/_plist/2018.json",
		"../patches/_plist/2019.json",
		"../patches/_plist/2020.json",
		"../patches/_plist/2021.json",
		"../patches/_plist/2022.json",
		"../patches/_plist/2023.json",
		"../patches/_plist/2024.json",
	}

	for _, file := range files {
		err := loadFromJson(file)
		if err != nil {
			panic(err)
		}
	}
}
