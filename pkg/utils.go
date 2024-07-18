package router

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
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

func displayError(c *gin.Context, message string, err error) {
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

func zipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func getFilesInDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
