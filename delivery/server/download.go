package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//getSingleFileHandler

func getSingleFileHandler(ctx *gin.Context) {

	uri := ctx.Param("fileName")
	if uri == "" {
		ctx.String(http.StatusInternalServerError, "uri is required")
		return
	}

	ctx.File(filepath.Join("./data/files", uri))

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
