package thumbnail

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreaskoch/go-fswatch"
	"github.com/disintegration/imaging"
)

func watchForChanges(originalsDir, thumbnailsDir string) {
	recurse := true // include all sub directories

	skipDotFilesAndFolders := func(path string) bool {
		return strings.HasPrefix(filepath.Base(path), ".")
	}

	checkIntervalInSeconds := 1

	folderWatcher := fswatch.NewFolderWatcher(
		originalsDir,
		recurse,
		skipDotFilesAndFolders,
		checkIntervalInSeconds,
	)

	folderWatcher.Start()
	for folderWatcher.IsRunning() {

		select {
		case changes := <-folderWatcher.ChangeDetails():
			for _, file := range changes.New() {
				thumbnailPath := thumbnailsDir + "/" + filepath.Base(file)
				makeThumbnail(file, thumbnailPath)
			}
			for _, file := range changes.Modified() {
				thumbnailPath := thumbnailsDir + "/" + filepath.Base(file)
				makeThumbnail(file, thumbnailPath)
			}
			for _, file := range changes.Moved() {
				thumbnailPath := thumbnailsDir + "/" + filepath.Base(file)
				_ = os.Remove(thumbnailPath)
			}
		}
	}
}

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
				makeThumbnail(originalPicturePath, thumbnailPath)
			}
		}
	}

	watchForChanges(originalsDir, thumbnailsDir)
}

func makeThumbnail(originalPicturePath, thumbnailPath string) {
	src, err := imaging.Open(originalPicturePath, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatalf("failed to open original: %v", err)
	}

	thumbnail := imaging.Thumbnail(src, 200, 200, imaging.Box)

	err = imaging.Save(thumbnail, thumbnailPath, imaging.JPEGQuality(80))
	if err != nil {
		log.Fatalf("failed to save thumbnail: %v", err)
	}

	log.Print(".")
}
