package server

import (
	"path"

	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
	"github.com/guilherme-gm/ro-vis/extractor/internal/database/repository"
)

type ServerType string

const (
	ServerTypeKROMain ServerType = "kRO-main"
	ServerTypeLATAM   ServerType = "latam"
	ServerTypeTest    ServerType = "test"
)

type Server struct {
	Type             ServerType
	PatchListUrl     string
	PatchFolderUrl   string
	LocalPatchFolder string
	DatabaseName     string
	Database         *database.Database
	Repositories     *repository.Repository
}

func New(svType ServerType, patchListUrl string, patchFolderUrl string, localPatchFolder string, databaseName string) *Server {
	db := database.NewDatabase(databaseName)

	return &Server{
		Type:             svType,
		PatchListUrl:     patchListUrl,
		PatchFolderUrl:   patchFolderUrl,
		LocalPatchFolder: localPatchFolder,
		DatabaseName:     databaseName,
		Database:         db,
		Repositories:     repository.NewRepository(db),
	}
}

func (s *Server) GetPatchesFolder() string {
	return path.Join("..", "patches", s.LocalPatchFolder)
}

func (s *Server) GetRawPatchesFolder() string {
	return path.Join("..", "patches", s.LocalPatchFolder, "_raw")
}

func (s *Server) GetPatchFile(patchName string) string {
	return path.Join(s.GetRawPatchesFolder(), patchName)
}

func (s *Server) GetExtractedPatchFolder(patchName string) string {
	return path.Join(s.GetPatchesFolder(), patchName)
}

var kroMain *Server
var latam *Server

func GetKROMain() *Server {
	if kroMain == nil {
		kroMain = New(
			ServerTypeKROMain,
			"", // @TODO: patch list url
			"", // @TODO: patch folder url
			"kro",
			"rovis-kro")
	}
	return kroMain
}

func GetLATAM() *Server {
	if latam == nil {
		latam = New(
			ServerTypeLATAM,
			"https://ro1patch.gnjoylatam.com/LIVE/patchinfo/patch.txt",
			"https://ro1patch.gnjoylatam.com/LIVE/patchfile/",
			"latam",
			"rovis-latam")
	}
	return latam
}

func GetTestServer() *Server {
	return &Server{
		Type:             ServerTypeTest,
		PatchListUrl:     "", // @TODO: patch list url
		PatchFolderUrl:   "", // @TODO: patch folder url
		LocalPatchFolder: "test",
		DatabaseName:     "rovis-test",
		Database:         nil,
		Repositories:     repository.NewRepository(nil),
	}
}

func GetServers() []*Server {
	return []*Server{
		// GetTestServer(), // not used for production
		GetKROMain(),
		GetLATAM(),
	}
}
