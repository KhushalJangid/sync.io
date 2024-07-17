package router

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"github.com/skip2/go-qrcode"
)

func traceDevices(c *gin.Context, connected_devices map[string]map[string]string) {
	remoteIp := c.RemoteIP()
	val, ok := connected_devices[remoteIp]
	currentTime := time.Now()
	if ok {
		val["last_connected"] = currentTime.Format("15:04:00 PM")
	} else {
		ua := useragent.Parse(c.GetHeader("user-agent"))
		fmt.Println("New connection")
		fmt.Println("Name:", ua.Name, "v", ua.Version)
		fmt.Println("OS:", ua.OS, "v", ua.OSVersion)
		fmt.Println("Device:", ua.Device)
		var name string
		if len(connected_devices) == 0 {
			name = ua.Name + " (Host)"
		} else {
			name = ua.Name
		}
		connected_devices[remoteIp] = map[string]string{
			"name":            name,
			"os":              ua.OS,
			"ip":              remoteIp,
			"first_connected": currentTime.Format("15:04:00 PM"),
			"last_connected":  currentTime.Format("15:04:00 PM"),
		}
	}
}

func generateQR() string {
	url := "http://" + GetOutboundIP().String() + ":8080"
	qrCode, _ := qrcode.New(url, qrcode.High)
	bytes, err := qrCode.PNG(256)
	if err != nil {
		panic(err)
	}
	imgBase64Str := base64.StdEncoding.EncodeToString(bytes)
	return imgBase64Str
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func ListFiles(ctx *gin.Context) {
	path := getDir()
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	output := []map[string]string{}
	for _, e := range entries {
		info, _ := e.Info()
		var size string
		if info.Size() < 1024 {
			size = strconv.Itoa(int(info.Size())) + " bytes"
		} else if 1024 < info.Size() && info.Size() < 1048576 {
			size = strconv.FormatFloat(float64(info.Size())/float64(1024), 'f', 2, 64) + " Kbs"
		} else {
			size = strconv.FormatFloat(float64(info.Size())/float64(1048576), 'f', 2, 64) + " Mbs"
		}
		output = append(output, map[string]string{
			"name": e.Name(),
			"size": size,
		})

	}
	ctx.HTML(http.StatusOK, "file_list.html", gin.H{"files": output})
}

func UploadFiles(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["file"]
	// path, err := filepath.Abs("cache")
	// if err != nil {
	// 	DisplayError(ctx, "Error in Uploading Multiple file , keep the key = 'file' to avoid this error ", err)
	// }
	path := getDir()
	for _, file := range files {
		log.Println(file.Filename)
		err := ctx.SaveUploadedFile(file, filepath.Join(path, file.Filename))
		if err != nil {
			DisplayError(ctx, "Error in Uploading Multiple file ,as the storage dir. not found ", err)
		}
	}

	fmt.Println(ctx.PostForm("key"))
	// ctx.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	ctx.Redirect(http.StatusFound, "/files")

}

func DownloadFile(ctx *gin.Context) {
	path, err := filepath.Abs("cache")
	if err != nil {
		log.Fatal(err)
	}
	fileName := ctx.Param("filename")

	ctx.FileAttachment(filepath.Join(path, fileName), fileName)
}
func PreviewFile(ctx *gin.Context) {
	path, err := filepath.Abs("cache")
	if err != nil {
		log.Fatal(err)
	}
	fileName := ctx.Param("filename")

	ctx.File(filepath.Join(path, fileName))
}
func DeleteFile(ctx *gin.Context) {
	path, err := filepath.Abs("cache")
	if err != nil {
		log.Fatal(err)
	}
	fileName := ctx.Param("filename")
	e := os.Remove(filepath.Join(path, fileName))
	if e == nil {
		ctx.Redirect(http.StatusFound, "/files")
	} else {
		panic(e)
	}

}

func DisplayError(c *gin.Context, message string, err error) {
	print(message + " <<<- this error @ this endpoint ->>> ") // to print all erroe at console
	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Note": message, "Error": err})
}

func getDir() string {
	path, err := filepath.Abs("cache")
	if err != nil {
		panic(err)
	}
	exist, _ := exists(path)
	if exist {
		return path
	} else {
		os.MkdirAll(path, 0777)
		return path
	}
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
