package thumbnail

import (
	"log"
	"io/ioutil"
	"sort"
)

func RetrieveAll(thumbnailsDir string) []string {
	thumbnailFiles, err := ioutil.ReadDir(thumbnailsDir)
	if err != nil {
		log.Fatal(err)
	}

	var thumbnailNames []string
	for _, thumbnailFile := range thumbnailFiles {
		thumbnailNames = append(thumbnailNames, thumbnailFile.Name())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(thumbnailNames)))

	return thumbnailNames
}
