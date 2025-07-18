package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchDownloader"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/grf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/rgz"
)

const notFoundString = "<?xml version='1.0' encoding='UTF-8'?><Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>"

func downloadPatches(server *server.Server) {
	fmt.Println("-- Checking for new updates")

	latest, err := server.Repositories.PatchRepository.GetLatestPatch(nil)
	if err != nil {
		panic(err)
	}

	patchServer, err := patchDownloader.NewPatchServer(server.PatchListUrl, server.PatchFolderUrl)
	if err != nil {
		panic(err)
	}

	for _, patch := range patchServer.PatchList {
		if patch.Id <= int(latest.Id) {
			continue
		}

		filePath := path.Join("..", "patches", server.LocalPatchFolder, "_raw", patch.Name)
		err := patchServer.DownloadPatch(&patch, filePath)
		if err != nil {
			panic(err)
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		newPatch := domain.Patch{
			Id:     int32(patch.Id),
			Name:   patch.Name,
			Date:   domain.TryGetPatchDate(patch.Name),
			Files:  []string{},
			Status: domain.PatchStatusPending,
		}
		if newPatch.Date.IsZero() {
			fmt.Printf("Unknown patch date: %s\n", patch.Name)
			fmt.Printf("Stopping patch download for %s\n", server.Type)
			return
		}

		if patch.Disabled {
			newPatch.Status = domain.PatchStatusSkipped
			fmt.Printf("Patch %s is disabled\n", patch.Name)
		}

		if string(fileContent) == notFoundString {
			if !patch.Disabled {
				fmt.Printf("Patch %s not found, but expected\n", patch.Name)
				return
			}

			newPatch.Status = domain.PatchStatusGone
			fmt.Printf("Patch %s not found, but it was disabled already\n", patch.Name)
			server.Repositories.PatchRepository.InsertPatch(nil, &newPatch)

			continue
		}

		if strings.HasSuffix(strings.ToLower(patch.Name), ".rgz") {
			rgzFile, err := rgz.Open(filePath)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Files in %s:\n", patch.Name)
			for _, entry := range rgzFile.Entries {
				fmt.Printf("%s\n", entry.Name)

				if entry.EntryType == rgz.EntryType_File {
					newPatch.Files = append(newPatch.Files, entry.Name)
				}
			}
		} else if strings.HasSuffix(strings.ToLower(patch.Name), ".gpf") {
			gpfFile, err := grf.Open(filePath)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Files in %s:\n", patch.Name)
			for _, file := range gpfFile.FileTable.Files {
				fmt.Printf("%s\n", file.FileName)

				if file.Flags == grf.EntryType_File {
					newPatch.Files = append(newPatch.Files, file.FileName)
				}
			}
		} else {
			// @TODO: Add some kind of alert for this
			fmt.Printf("Unsupported patch format: %s\n", patch.Name)
			fmt.Printf("Stopping patch download for %s\n", server.Type)
			return
		}

		server.Repositories.PatchRepository.InsertPatch(nil, &newPatch)
	}
}

func main() {
	fmt.Println("RO Vis extractor - Miner")
	conf.LoadExtractor()

	// @TODO: server.GetServers()
	for _, server := range []*server.Server{server.GetLATAM()} {
		fmt.Println("------ Mining " + server.DatabaseName + " ------")
		downloadPatches(server)
	}

	fmt.Println("Success")
}
