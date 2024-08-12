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
	_, ok := connected_devices[remoteIp]
	currentTime := time.Now()
	if ok {
		// val["last_connected"] = StringBool{Str: currentTime.Format("15:04:00 PM")}
		return
	} else {
		ua := useragent.Parse(c.GetHeader("user-agent"))
		print(ua.Device)
		if len(connected_devices) == 0 {
			device_map := make(map[string]StringBool)

			device_map["name"] = StringBool{Str: ua.Name + " (Host)"}
			device_map["os"] = StringBool{Str: ua.OS}
			device_map["mobile"] = StringBool{Flag: ua.Mobile}
			device_map["ip"] = StringBool{Str: remoteIp}
			device_map["connected"] = StringBool{Str: currentTime.Format("15:04:00 PM")}
			// device_map["last_connected"] = StringBool{Str: currentTime.Format("15:04:00 PM")}
			connected_devices[remoteIp] = device_map
			c.Redirect(302, "/setpassword")
		} else {
			device_map := make(map[string]StringBool)

			device_map["name"] = StringBool{Str: ua.Name}
			device_map["os"] = StringBool{Str: ua.OS}
			device_map["mobile"] = StringBool{Flag: ua.Mobile}
			device_map["ip"] = StringBool{Str: remoteIp}
			device_map["connected"] = StringBool{Str: currentTime.Format("15:04:00 PM")}
			// device_map["last_connected"] = StringBool{Str: currentTime.Format("15:04:00 PM")}

			connected_devices[remoteIp] = device_map
			if channel.password != "" {
				c.Redirect(302, "/verifypassword")
			}
		}
	}
}

func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}
}
