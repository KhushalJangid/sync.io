package router

import (
	"embed"
	"html/template"
	io "io/fs"
	"net/http"

	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//go:embed templates/tailwind/*
var templateFS embed.FS

//go:embed templates/static/*
var staticFiles embed.FS
var channel *Channel

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	tmpl := template.Must(template.ParseFS(templateFS, "templates/tailwind/*"))
	channel = new(Channel)
	channel.connected_devices = make(map[string]map[string]StringBool)
	r.SetHTMLTemplate(tmpl)
	r.Use(sessionMiddleware())
	r.Use(gin.Recovery())
	r.Use(brotli.Brotli(brotli.DefaultCompression))
	static := r.Group("/")
	{
		static.Use(cacheMiddleware())
		fs, _ := io.Sub(staticFiles, "templates/static")
		static.StaticFS("/static/", http.FS(fs))
	}
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/files", ListFiles)
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{"csrf": csrf.GetToken(c)})
	})
	r.POST("/upload", UploadFiles)
	r.GET("/connected_devices", func(c *gin.Context) {
		c.HTML(http.StatusOK, "connectedDevices.html", gin.H{"devices": channel.connected_devices})
	})
	r.GET("/qr", func(c *gin.Context) {
		c.HTML(http.StatusOK, "generateQR.html", gin.H{"qr": generateQR()})
	})
	r.GET("/download/:filename", DownloadFile)
	r.GET("/preview/:filename", PreviewFile)
	r.GET("/delete/:filename", DeleteFile)
	r.GET("/downloadAll", DownloadAllFiles)
	r.GET("/deleteAll", DeleteAllFiles)
	r.GET("/setpassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setPassword.html", gin.H{"csrf": csrf.GetToken(c)})
	})
	r.GET("/verifypassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "verifyPassword.html", gin.H{"csrf": csrf.GetToken(c)})
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
