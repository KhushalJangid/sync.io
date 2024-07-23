package router

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templateFS embed.FS
var channel *Channel

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.LoadHTMLGlob("templates/*")
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*"))
	channel = new(Channel)
	channel.connected_devices = map[string]map[string]string{}
	r.SetHTMLTemplate(tmpl)
	r.Use(sessionMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})
	r.GET("/files", ListFiles)
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})
	r.POST("/upload", UploadFiles)
	r.GET("/connected_devices", func(c *gin.Context) {
		c.HTML(http.StatusOK, "connected_devices.html", gin.H{"devices": channel.connected_devices})
	})
	r.GET("/qr", func(c *gin.Context) {
		c.HTML(http.StatusOK, "qr.html", gin.H{"qr": generateQR()})
	})
	r.GET("/download/:filename", DownloadFile)
	r.GET("/preview/:filename", PreviewFile)
	r.GET("/delete/:filename", DeleteFile)
	r.GET("/downloadAll", DownloadAllFiles)
	r.GET("/deleteAll", DeleteAllFiles)
	r.GET("/setpassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setpassword.html", gin.H{})
	})
	r.GET("/verifypassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "verifypassword.html", gin.H{})
	})
	r.POST("/setpassword", func(c *gin.Context) {
		password := c.PostForm("password")
		setPassword(password)
		c.Redirect(302, "/")
	})
	r.POST("/verifypassword", func(c *gin.Context) {
		password := c.PostForm("password")
		iscorrect := verifyPassword(password)
		if iscorrect {
			c.Redirect(302, "/")
		}
	})

	return r
}
