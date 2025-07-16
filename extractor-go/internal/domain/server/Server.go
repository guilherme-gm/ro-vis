package server

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/database"
)

type ServerType string

const (
	ServerTypeKROMain ServerType = "kRO-main"
	ServerTypeLATAM   ServerType = "latam"
)

type Server struct {
	Type             ServerType
	PatchListUrl     string
	PatchFolderUrl   string
	LocalPatchFolder string
	DatabaseName     string
	Database         *database.Database
}

func New(svType ServerType, patchListUrl string, patchFolderUrl string, localPatchFolder string, databaseName string) *Server {
	return &Server{
		Type:             svType,
		PatchListUrl:     patchListUrl,
		PatchFolderUrl:   patchFolderUrl,
		LocalPatchFolder: localPatchFolder,
		DatabaseName:     databaseName,
		Database:         database.NewDatabase(databaseName),
	}
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
			"ro-vis-kro-main")
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
			"ro-vis-latam")
	}
	return latam
}
