package server

import (
	"context"
	"github.com/gin-gonic/gin"
	proto "imageCache/grpc/gen/proto/imageCache/v1"
	"net/http"
	"path/filepath"
)

func (s *ServerGRPC) ListFiles(context.Context, *proto.Null) (*proto.ListFilesResponse, error) {

	var result = proto.ListFilesResponse{}

	res, err := listAllFiles()
	if err != nil {
		return nil, err
	}

	for _, file := range res {
		result.Files = append(result.Files, &proto.File{
			FileName: filepath.Base(file),
			Uri:      file,
		})
	}
	return &result, nil
}

func listFilesHandler(ctx *gin.Context) {
	response, er := listAllFiles()
	if er != nil {
		ctx.String(http.StatusInternalServerError, er.Error())
		return
	}

	ctx.JSONP(http.StatusOK, response)
}

func listAllFiles() ([]string, error) {
	return filepath.Glob("./data/files/*")
}
