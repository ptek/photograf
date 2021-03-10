package thumbnail

import (
	"fmt"
	"github.com/edwvee/exiffix"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Generate(originalsDir, thumbnailsDir string, regenerate bool) {
	files, err := ioutil.ReadDir(originalsDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(thumbnailsDir, 0755); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		originalPicturePath := originalsDir + "/" + file.Name()

		if strings.Contains(strings.ToLower(originalPicturePath), ".jpg") || strings.Contains(strings.ToLower(originalPicturePath), ".jpeg") {
			thumbnailPath := thumbnailsDir + "/" + file.Name()

			if _, err := os.Stat(thumbnailPath); !(err == nil && regenerate == false) {
				resizeAndCrop(originalPicturePath, thumbnailPath)
			}
		}
	}
}

func resizeAndCrop(originalPicturePath, thumbnailPath string) {
	file, err := os.Open(originalPicturePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	picture, _, err := exiffix.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	pictureDimensions := picture.Bounds()

	var resized image.Image
	if pictureDimensions.Dx() > pictureDimensions.Dy() {
		resized = resize.Resize(0, 200, picture, resize.NearestNeighbor)
	} else {
		resized = resize.Resize(200, 0, picture, resize.NearestNeighbor)
	}

	outputFile, err := os.Create(thumbnailPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	jpeg.Encode(outputFile, resized, nil)
	fmt.Print(".")
}
