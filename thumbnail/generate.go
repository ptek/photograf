package thumbnail

import (
  "fmt"
  "image"
  "image/jpeg"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "strings"

  "github.com/andreaskoch/go-fswatch"
  "github.com/edwvee/exiffix"
  "github.com/nfnt/resize"
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
