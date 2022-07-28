package thumbnail

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fsnotify/fsnotify"
)

func watchForChanges(originalsDir, thumbnailsDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					thumbnailPath := thumbnailsDir + "/" + filepath.Base(event.Name)
					makeThumbnail(event.Name, thumbnailPath)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					thumbnailPath := thumbnailsDir + "/" + filepath.Base(event.Name)
					makeThumbnail(event.Name, thumbnailPath)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					thumbnailPath := thumbnailsDir + "/" + filepath.Base(event.Name)
					makeThumbnail(event.Name, thumbnailPath)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatal("error:", err)
			}
		}
	}()

	err = watcher.Add(originalsDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
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

	log.Printf("Created thumbnail for: %s", originalPicturePath)
}
