package testfiles

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "..")

	// Test files folder (this module folder)
	TestFiles = filepath.Dir(b)
)

func GetFilePath(fileName string) string {
	return filepath.Join(TestFiles, fileName)
}
