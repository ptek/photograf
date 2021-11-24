package thumbnail

import (
	"io/ioutil"
	"log"
	"sort"
)

func RetrieveGalleryItems(thumbnailsDir string) []string {
	thumbnailFiles, err := ioutil.ReadDir(thumbnailsDir)
	if err != nil {
		log.Fatal(err)
	}

	var galleryItems []string
	for _, item := range thumbnailFiles {
		galleryItems = append(galleryItems, item.Name())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(galleryItems)))

	return galleryItems
}
