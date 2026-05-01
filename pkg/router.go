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
	r.Use(gin.Recovery())
	r.Use(brotli.Brotli(brotli.DefaultCompression))

	// Static assets — no auth required (needed by login pages too).
	static := r.Group("/")
	{
		static.Use(cacheMiddleware())
		fs, _ := io.Sub(staticFiles, "templates/static")
		static.StaticFS("/static/", http.FS(fs))
	}

	// Auth pages — intentionally outside the protected group.
	r.GET("/setpassword", func(c *gin.Context) {
		channel.mu.RLock()
		deviceCount := len(channel.connected_devices)
		channel.mu.RUnlock()
		if deviceCount > 0 {
			c.Redirect(http.StatusFound, "/verifypassword")
			return
		}
		c.HTML(http.StatusOK, "setPassword.html", gin.H{"csrf": csrf.GetToken(c)})
	})
	r.POST("/setpassword", func(c *gin.Context) {
		channel.mu.RLock()
		deviceCount := len(channel.connected_devices)
		channel.mu.RUnlock()
		if deviceCount > 0 {
			c.Redirect(http.StatusFound, "/verifypassword")
			return
		}
		setPassword(c.PostForm("password"))
		issueToken(c, true)
		c.Redirect(http.StatusFound, "/")
	})
	r.GET("/verifypassword", func(c *gin.Context) {
		session := sessions.Default(c)
		token, _ := session.Get("auth_token").(string)
		channel.mu.RLock()
		_, valid := channel.connected_devices[token]
		channel.mu.RUnlock()
		if valid {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "verifyPassword.html", gin.H{"csrf": csrf.GetToken(c)})
	})
	r.POST("/verifypassword", func(c *gin.Context) {
		if verifyPassword(c.PostForm("password")) {
			issueToken(c, false)
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusUnauthorized, "verifyPassword.html", gin.H{
			"csrf":  csrf.GetToken(c),
			"error": "Incorrect password. Please try again.",
		})
	})

	// Protected routes — all require a valid server-side token.
	protected := r.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
		protected.GET("/files", ListFiles)
		protected.GET("/upload", func(c *gin.Context) {
			c.HTML(http.StatusOK, "upload.html", gin.H{"csrf": csrf.GetToken(c)})
		})
		protected.POST("/upload", UploadFiles)
		protected.GET("/connected_devices", func(c *gin.Context) {
			session := sessions.Default(c)
			myToken, _ := session.Get("auth_token").(string)
			channel.mu.RLock()
			me := channel.connected_devices[myToken]
			devices := channel.connected_devices
			channel.mu.RUnlock()
			c.HTML(http.StatusOK, "connectedDevices.html", gin.H{
				"devices": devices,
				"isHost":  me["isHost"].Flag,
			})
		})
		protected.GET("/qr", func(c *gin.Context) {
			c.HTML(http.StatusOK, "generateQR.html", gin.H{"qr": generateQR()})
		})
		protected.GET("/download/:filename", DownloadFile)
		protected.GET("/preview/:filename", PreviewFile)
		protected.GET("/delete/:filename", DeleteFile)
		protected.GET("/downloadAll", DownloadAllFiles)
		protected.GET("/deleteAll", DeleteAllFiles)
		protected.GET("/revoke/:token", RevokeDevice)
	}

	return r
}
