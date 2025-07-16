package patchDownloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PatchItem struct {
	Id       int
	Name     string
	Disabled bool
}

type PatchServer struct {
	PatchListUrl   string
	PatchFolderUrl string
	PatchList      []PatchItem
}

func parsePatchList(list string) ([]PatchItem, error) {
	var patchList []PatchItem
	lines := strings.Split(list, "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		disabled := false
		if line == "" {
			continue
		}

		if line[0:2] == "//" {
			disabled = true
			line = strings.TrimSpace(line[2:])
		}

		fields := strings.Split(line, " ")
		patchId, err := strconv.Atoi(fields[0])
		if errors.Is(err, strconv.ErrSyntax) {
			fmt.Printf("Skipping line %d: %v is not a number.\n", idx, fields[0])
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse patch list line %d: %v is not a number. Error: %v", idx, fields[0], err)
		}

		fileName := strings.TrimSpace(fields[1])
		if fileName == "" {
			return nil, fmt.Errorf("failed to parse patch list line %d: %v is not a valid file name", idx, fields[1])
		}

		patchList = append(patchList, PatchItem{
			Id:       patchId,
			Name:     fileName,
			Disabled: disabled,
		})
	}

	return patchList, nil
}

func NewPatchServer(patchListUrl string, patchFolderUrl string) *PatchServer {
	resp, err := http.Get(patchListUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	patchList, err := parsePatchList(string(data))
	if err != nil {
		panic(err)
	}

	return &PatchServer{
		PatchListUrl:   patchListUrl,
		PatchFolderUrl: patchFolderUrl,
		PatchList:      patchList,
	}
}

func (ps *PatchServer) DownloadPatch(patch *PatchItem, savePath string) error {
	resp, err := http.Get(ps.PatchFolderUrl + patch.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	os.WriteFile(savePath, data, 0644)

	return nil
}
