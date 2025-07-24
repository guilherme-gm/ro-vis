package patchfile

type PatchFile interface {
	Extract(filePath string, rootFolder string) error
}
