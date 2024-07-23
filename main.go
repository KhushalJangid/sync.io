package main

import (
	"fmt"
	"time"

	router "Sync.io/pkg"
)

func main() {
	currentTime := time.Now()
	fmt.Println(currentTime.Format("Monday 02-Jan-2006 15:04:00 PM"))
	fmt.Println("Starting server...")
	fmt.Println("Listening on : ", "http://"+router.GetOutboundIP()+":8080/")
	r := router.Router()
	r.Run("0.0.0.0:8080")
}
