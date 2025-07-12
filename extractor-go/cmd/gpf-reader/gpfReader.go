/**
 * "Development" GPF Reader
 * Attempts to read a GPF file. for testing purposes
 */

package main

import (
	"fmt"
	"os"

	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/patchfile/rgz"
)

func main() {
	fmt.Println("RO Vis - GPF Reader")

	// gpfFile, err := grf.NewGpfFile("../2025-03-24_live_client_20_22_1742783771.gpf")
	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Printf("%+v\n", gpfFile.FileTable.Files[0])
	// // fmt.Printf("%+v\n", gpfFile.FileTable.Files[1])
	// // fmt.Printf("%+v\n", gpfFile.FileTable.Files[2])
	// fmt.Printf("%+v\n", gpfFile.FileTable.Files[3])
	// // fmt.Printf("%+v\n", gpfFile.FileTable.Files[4])
	// // fmt.Printf("%+v\n", gpfFile.FileTable.Files[5])

	// fmt.Printf("%+v\n", gpfFile.FileTable.Files[3])
	// if err := gpfFile.FileTable.Files[3].Extract("../2025-03-24_live_client_20_22_1742783771.gpf"); err != nil {
	// 	panic(err)
	// }

	rgzFile, err := rgz.Open("../2025-03-14_live_client_1_9_1741940954.rgz")
	if err != nil {
		panic(err)
	}

	os.WriteFile("test.lub", rgzFile.Entries[3].Data, 0644)

	fmt.Printf("%+v\n", rgzFile.Entries)

	fmt.Println("Success")
}
