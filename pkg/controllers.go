package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	path := getDir()
	for _, file := range files {
		log.Println(file.Filename)
		err := ctx.SaveUploadedFile(file, filepath.Join(path, file.Filename))
		if err != nil {
			displayError(ctx, "Error in Uploading Multiple file ,as the storage dir. not found ", err)
		}
	}

	fmt.Println(ctx.PostForm("key"))
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
func DownloadAllFiles(ctx *gin.Context) {
	path, _ := filepath.Abs("cache")
	files, err := getFilesInDir(path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read directory: %v", err)
		return
	}

	zipFile := "files.zip"
	if err := zipFiles(zipFile, files); err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to zip files: %v", err)
		return
	}

	ctx.FileAttachment(zipFile, "files.zip")
	defer os.Remove(zipFile)

}
func PreviewFile(ctx *gin.Context) {
	path, err := filepath.Abs("cache")
	if err != nil {
		log.Fatal(err)
	}
	fileName := ctx.Param("filename")

	ctx.File(filepath.Join(path, fileName))
}
func DeleteAllFiles(ctx *gin.Context) {
	path, err := filepath.Abs("cache")
	if err != nil {
		log.Fatal(err)
	}
	entries, _ := os.ReadDir(path)
	for _, e := range entries {
		os.Remove(filepath.Join(path, e.Name()))
	}
	ctx.Redirect(http.StatusFound, "/files")
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
