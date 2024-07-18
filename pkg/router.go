package router

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templateFS embed.FS

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.LoadHTMLGlob("templates/*")
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*"))
	r.SetHTMLTemplate(tmpl)
	connected_devices := map[string]map[string]string{}
	r.GET("/", func(c *gin.Context) {
		traceDevices(c, connected_devices)
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})
	r.GET("/files", ListFiles)
	r.GET("/upload", func(c *gin.Context) {
		traceDevices(c, connected_devices)
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})
	r.POST("/upload", UploadFiles)
	r.GET("/connected_devices", func(c *gin.Context) {
		traceDevices(c, connected_devices)
		c.HTML(http.StatusOK, "connected_devices.html", gin.H{"devices": connected_devices})
	})
	r.GET("/qr", func(c *gin.Context) {
		traceDevices(c, connected_devices)
		c.HTML(http.StatusOK, "qr.html", gin.H{"qr": generateQR()})
	})
	r.GET("/download/:filename", DownloadFile)
	r.GET("/preview/:filename", PreviewFile)
	r.GET("/delete/:filename", DeleteFile)
	r.GET("/downloadAll", DownloadAllFiles)
	r.GET("/deleteAll", DeleteAllFiles)
	return r
}
