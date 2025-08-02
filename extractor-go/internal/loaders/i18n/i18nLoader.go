package i18n

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type I18nLoader struct {
	repository *repository.I18nRepository
}

// GetRelevantFiles returns a list of all files that are relevant to this loader's parsers.
func (l *I18nLoader) GetRelevantFiles() []*regexp.Regexp {
	return []*regexp.Regexp{
		scFiles,
		scJsonFile,
	}
}

func NewI18nLoader(server *server.Server) *I18nLoader {
	return &I18nLoader{
		repository: server.Repositories.I18nRepository,
	}
}

func (l *I18nLoader) LoadPatch(tx *sql.Tx, basePath string, update domain.Update) {
	// Per my understanding, they throw all i18n files from the entire game into the grf already,
	// and are just enabling them via sc.json as needed.
	// Thus, the approach taken here is to always load everything, but mark them as active/inactive
	// based on the sc.json file.

	if !update.HasChangedAnyFiles(l.GetRelevantFiles()) {
		fmt.Println("Skipped - No meaningful file")
		return
	}

	fmt.Println("> Decoding...")

	// Get items currently loaded
	fmt.Println("> Fetching current list...")
	currentI18ns, err := l.repository.GetCurrentI18ns(tx)
	if err != nil {
		panic(err)
	}

	i18nMap := make(map[string]*domain.I18n)
	i18nFileIds := make(map[string][]string)
	for _, q := range *currentI18ns {
		i18nMap[q.I18nId] = &q
		i18nFileIds[q.ContainerFile] = append(i18nFileIds[q.ContainerFile], q.I18nId)
	}

	// Define list of csv's that matter (from json or from previous version)
	fmt.Println("> Fetching active file list...")
	activeFiles := make(map[string]bool)
	fileListChange, err := update.GetChangeForFile("data/i18n/sc/sc.json")
	if err == nil {
		// there is a new file list
		data, err := os.ReadFile(basePath + "/" + fileListChange.Patch + "/data/i18n/sc/sc.json")
		if err != nil {
			panic(err)
		}

		var fileListArr []string
		err = json.Unmarshal(data, &fileListArr)
		if err != nil {
			panic(err)
		}

		for _, v := range fileListArr {
			activeFiles[fmt.Sprintf("data/i18n/sc/%s.csv", v)] = true
		}
	} else if errors.Is(err, domain.NewNotFoundError("")) {
		// there is no new file list
		for _, q := range *currentI18ns {
			if q.Active {
				activeFiles[q.ContainerFile] = true
			}
		}
	} else {
		panic(err)
	}

	// Find CSVs to load
	fmt.Println("> Fetching changes...")
	changes, err := update.GetChangesForFile(scFiles)
	if err != nil {
		panic(err)
	}

	var newI18ns []domain.I18n
	var updatedI18ns []domain.I18n
	var idsToDelete []string

	// Load files
	// @TODO: new records are not being activated. And files that got activated in the update is not being marked as such.
	fmt.Println("> Loading files...")
	targetParser := NewI18nV1Parser()
	for _, change := range changes {
		fileEntries := targetParser.Parse(basePath, &change)
		fileIds := make(map[string]bool)
		fileIdsToDelete := make(map[string]bool)

		for _, entry := range i18nFileIds[change.File] {
			fileIdsToDelete[entry] = true
		}

		for _, fileEntry := range fileEntries {
			fileIds[fileEntry.I18nId] = true
			delete(fileIdsToDelete, fileEntry.I18nId)

			existingEntry := i18nMap[fileEntry.I18nId]
			if existingEntry == nil {
				newI18ns = append(newI18ns, fileEntry)
				continue
			}

			if !existingEntry.Equals(fileEntry) {
				fileEntry.PreviousHistoryID = existingEntry.HistoryID
				updatedI18ns = append(updatedI18ns, fileEntry)
			}
		}

		for deletedId := range fileIdsToDelete {
			idsToDelete = append(idsToDelete, deletedId)
		}
	}

	// Inactivate files that are not used
	fmt.Println("> Inactivating files...")
	for fileName, ids := range i18nFileIds {
		if _, ok := activeFiles[fileName]; !ok {
			for _, id := range ids {
				existingEntry := i18nMap[id]
				if existingEntry == nil {
					continue
				}

				newEntry := *existingEntry
				newEntry.HistoryID = domain.NullableInt64{Valid: false}
				newEntry.PreviousHistoryID = existingEntry.HistoryID
				newEntry.Active = false
				updatedI18ns = append(updatedI18ns, newEntry)
			}
		}
	}

	// Saving
	fmt.Printf("> Saving new records... (%d records to save)\n", len(newI18ns))
	err = l.repository.AddI18nsToHistory(tx, update.Name(), &newI18ns)
	if err != nil {
		fmt.Printf("Error saving new records: %v\n", err)
		panic(err)
	}

	fmt.Printf("> Updating records... (%d records to update)\n", len(updatedI18ns))
	err = l.repository.AddI18nsToHistory(tx, update.Name(), &updatedI18ns)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Deleting records... (%d records to delete)\n", len(idsToDelete))
	for _, deletedId := range idsToDelete {
		err := l.repository.AddDeletedI18n(tx, update.Name(), i18nMap[deletedId])
		if err != nil {
			panic(err)
		}
	}
}

func (l *I18nLoader) Name() string {
	return "i18n"
}
