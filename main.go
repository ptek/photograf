package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ptek/photograf/thumbnail"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed ui/*
var uiContentFS embed.FS

func main() {
	router := gin.Default()
	LoadHTMLFromEmbedFS(router, templatesFS, "templates/*")

	router.StaticFS("/ui", UiFS())

	router.Static("/originals", originalsDir())
	router.Static("/thumbnails", thumbnailsDir())

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.tmpl", gin.H{
			"status": "success",
		})
	})

	router.GET("/pictures/:offset", func(c *gin.Context) {
		offset_p := c.Param("offset")
		if offset_p == "" {
			offset_p = "0"
		}

		offset, _ := strconv.ParseInt(offset_p, 10, 64)
		if offset < 0 {
			log.Fatal("Offset must be positive.")
		}

		galleryItems := thumbnail.RetrieveGalleryItems(thumbnailsDir())

		if len(galleryItems) <= int(offset) {
			log.Println("All images have been loaded")
		} else {
			next_offset := offset + 48

			if len(galleryItems) < int(next_offset) {
				next_offset = int64(len(galleryItems))
			}

			items := galleryItems[offset:next_offset]
			c.HTML(http.StatusOK, "pictures.tmpl", gin.H{
				"thumbnails":  items,
				"next_offset": next_offset,
			})
		}
	})

	router.GET("/original/:item", func(c *gin.Context) {
		item := c.Param("item")

		items := thumbnail.RetrieveGalleryItems(thumbnailsDir())

		var nextItem = ""
		var previousItem = ""
		for i, itemUnderInspection := range items {
			if i+1 < len(items) {
				nextItem = items[i+1]
			} else {
				nextItem = items[0]
			}

			if i > 0 {
				previousItem = items[i-1]
			} else {
				previousItem = items[len(items)-1]
			}

			if itemUnderInspection == item {
				break
			}
		}

		c.HTML(http.StatusOK, "original.tmpl", gin.H{
			"item":         item,
			"nextItem":     nextItem,
			"previousItem": previousItem,
		})
	})

	router.GET("/empty", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	log.Println("http://127.0.0.1:3000")
	router.Run()
}

// TEMPLATES

func LoadHTMLFromEmbedFS(engine *gin.Engine, embedFS embed.FS, pattern string) {
	templ := template.Must(template.ParseFS(embedFS, pattern))
	log.Println("Templ: ", templ)
	engine.SetHTMLTemplate(templ)
}

// SETTINGS

func originalsDir() string {
	var originalsDir = os.Getenv("ORIGINALS")
	if len(originalsDir) == 0 {
		originalsDir = "./pictures"
	}
	return originalsDir
}

func thumbnailsDir() string {
	var thumbnailsDir = os.Getenv("THUMBNAILS")
	if len(thumbnailsDir) == 0 {
		thumbnailsDir = "./thumbnails"
	}
	return thumbnailsDir
}

func regenerate() bool {
	var regenerate, err = strconv.ParseBool(os.Getenv("REGENERATE"))
	if err != nil {
		regenerate = false
	}
	return regenerate
}

func certs() string {
	var dir = os.Getenv("CERTIFICATES")
	if len(dir) == 0 {
		dir = "./cert"
	}
	return dir
}

// UTILITY FUNCTIONS

func UiFS() http.FileSystem {
	fsys, err := fs.Sub(uiContentFS, "ui")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

func init() {
	go thumbnail.Generate(originalsDir(), thumbnailsDir(), regenerate())
}
