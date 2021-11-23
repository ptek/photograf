package main

import (
  "embed"
  "html/template"
  "log"
  "net/http"
  "os"
  "strconv"

  "github.com/gin-gonic/gin"

  "github.com/ptek/photograf/thumbnail"
)

var html = template.Must(template.New("https").Parse(`
<html>
  <head>
    <title>Photograf</title>
    <link rel="stylesheet" href="/ui/style.css">
  </head>
  <body>
    <div id="original_view"></div>
    <div id="gallery"></div>
  </body>
  <script src="/ui/index.js"></script>
  <script>photograf_main()</script>
</html>
`))

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

//go:embed ui/*
var uiContent embed.FS

func StaticHandler(c *gin.Context) {
  filePath := c.Param("filepath")
  data, err := uiContent.ReadFile("ui" + filePath)
  if err != nil {
    return
  }
  c.Writer.Write(data)
}

func init() {
  go thumbnail.Generate(originalsDir(), thumbnailsDir(), regenerate())
}

func main() {
  router := gin.Default()
  router.SetHTMLTemplate(html)

  router.GET("/ui/*filepath", StaticHandler)
  router.Static("/originals", originalsDir())
  router.Static("/thumbnails", thumbnailsDir())

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "https", gin.H{
      "status": "success",
    })
  })

  router.GET("/api/thumbnail/:picture_id", func(c *gin.Context) {
    thumbnailId := c.Param("picture_id")
    c.JSON(http.StatusOK, gin.H{"hello": thumbnailId})
  })

  router.GET("/api/thumbnails/:from/:to", func(c *gin.Context) {

    from, _ := strconv.ParseInt(c.Param("from"), 10, 64)
    to, _ := strconv.ParseInt(c.Param("to"), 10, 64)

    if from < 0 {
      log.Fatal("Field from should be positive.")
    }

    allThumbnails := thumbnail.RetrieveAll(thumbnailsDir())

    if len(allThumbnails) < int(to) {
      to = int64(len(allThumbnails))
    }
    if to < from {
      from = to
    }
    thumbnails := allThumbnails[from:to]

    c.JSON(http.StatusOK, gin.H{
      "thumbnails": thumbnails,
    })
  })

  router.Run()
}
