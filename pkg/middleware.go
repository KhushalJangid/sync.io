package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

func sessionMiddleware() gin.HandlerFunc {
	return traceDevices
}

func traceDevices(c *gin.Context) {
	remoteIp := c.RemoteIP()
	connected_devices := channel.connected_devices
	val, ok := connected_devices[remoteIp]
	currentTime := time.Now()
	if ok {
		val["last_connected"] = currentTime.Format("15:04:00 PM")
	} else {
		ua := useragent.Parse(c.GetHeader("user-agent"))
		if len(connected_devices) == 0 {
			connected_devices[remoteIp] = map[string]string{
				"name":            ua.Name + " (Host)",
				"os":              ua.OS,
				"ip":              remoteIp,
				"first_connected": currentTime.Format("15:04:00 PM"),
				"last_connected":  currentTime.Format("15:04:00 PM"),
			}
			c.Redirect(302, "/setpassword")
		} else {
			connected_devices[remoteIp] = map[string]string{
				"name":            ua.Name,
				"os":              ua.OS,
				"ip":              remoteIp,
				"first_connected": currentTime.Format("15:04:00 PM"),
				"last_connected":  currentTime.Format("15:04:00 PM"),
			}
			if channel.password != "" {
				c.Redirect(302, "/verifypassword")
			}
		}
	}
}
