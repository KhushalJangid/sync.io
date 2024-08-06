package main

import (
	"fmt"
	"time"

	router "Sync.io/pkg"
)

func main() {
	currentTime := time.Now()
	url := "http://" + router.GetOutboundIP() + ":8080/"
	fmt.Println(currentTime.Format("Monday 02-Jan-2006 15:04:00 PM"))
	fmt.Println("Starting server...")
	fmt.Println("Listening on : ", url)
	r := router.Router()
	router.OpenBrowser(url)
	r.Run("0.0.0.0:8080")
}
