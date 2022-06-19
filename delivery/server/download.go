package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imageCache/data/files"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//DownloadFileHandler

func DownloadFileHandler(ctx *gin.Context) {

	uri := ctx.Param("fileName")
	if uri == "" {
		ctx.String(http.StatusInternalServerError, "fileName is required e.g <server addr>/data/files/sample1.png")
		return
	}

	ctx.File(filepath.Join(files.GetLocation(), uri))

}

type file struct {
	FileName string `json:"file_name"`
	Data     []byte `json:"data"`
}

func getFileByFileName(fileName string) (*file, error) {

	data, err := ioutil.ReadFile("./data/files/" + fileName)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	return &file{
		FileName: fileName,
		Data:     data,
	}, nil
}
