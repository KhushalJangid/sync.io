package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token, _ := session.Get("auth_token").(string)

		channel.mu.RLock()
		_, valid := channel.connected_devices[token]
		deviceCount := len(channel.connected_devices)
		noPassword := channel.password == ""
		channel.mu.RUnlock()

		if valid {
			c.Next()
			return
		}

		if deviceCount == 0 {
			c.Redirect(http.StatusFound, "/setpassword")
			c.Abort()
			return
		}

		// No password set: auto-issue a token so the user gets in seamlessly.
		if noPassword {
			issueToken(c, false)
			c.Next()
			return
		}

		c.Redirect(http.StatusFound, "/verifypassword")
		c.Abort()
	}
}

// issueToken creates a new auth token for the device, stores it server-side,
// and writes it to the session cookie.
func issueToken(c *gin.Context, isHost bool) {
	ua := useragent.Parse(c.GetHeader("user-agent"))
	token := generateToken()
	name := ua.Name
	if isHost {
		name += " (Host)"
	}
	device := map[string]StringBool{
		"name":      {Str: name},
		"os":        {Str: ua.OS},
		"mobile":    {Flag: ua.Mobile},
		"ip":        {Str: c.RemoteIP()},
		"connected": {Str: time.Now().Format("15:04:00 PM")},
		"isHost":    {Flag: isHost},
		"token":     {Str: token},
	}

	channel.mu.Lock()
	channel.connected_devices[token] = device
	channel.mu.Unlock()

	session := sessions.Default(c)
	session.Set("auth_token", token)
	session.Save()
}

func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}
}
