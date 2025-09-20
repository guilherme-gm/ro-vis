package server

import (
	"strings"
)

var skipFiles = map[ServerType]map[string][]string{
	ServerTypeKROMain: {},
	ServerTypeLATAM: {
		"2025-08-27_live_client_1010_1010_1756282119.gpf": {
			// These files seems to be broken, as we can't decompress them properly (GRF Editor also fails)
			// "zlib: Unexpected EOF"
			"data/i18n/sc/a104a179a7789db1b0c7c40e9bd45ab137c2e9f9e639d2f28b5033a65c02dece.csv",
			"data/i18n/sc/fed3811f813f6a737a213a57464f694ad9a61fa41dce38b31eeb0d1d5ebf6214.csv",
		},
		"2025-09-12_live_client_1058_1059_1757645955.gpf": {
			// This file failed checksum check on 2025-09-12, it was the only patch for it, and the file looks ok
			// but implementing a bypass for the checksum here is too much work, it is a single sentence,
			// so we can leave without it.
			"data/i18n/sc/fed3811f813f6a737a213a57464f694ad9a61fa41dce38b31eeb0d1d5ebf6214.csv",
		},
	},
}

func ShouldSkipFile(server *Server, patch string, file string) bool {
	serverSkipMap, ok := skipFiles[server.Type]
	if !ok {
		return false
	}

	patchSkipList, ok := serverSkipMap[patch]
	if !ok {
		return false
	}

	for _, skipFile := range patchSkipList {
		if strings.EqualFold(skipFile, file) {
			return true
		}
	}

	return false
}
