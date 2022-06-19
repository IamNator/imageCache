package server

import (
	"fmt"
	proto "imageCache/grpc/gen/proto/imageCache/v1"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//writeToFp takes in a file pointer and byte array and writes the byte array into the file
//returns error if pointer is nil or error in writing to file
func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {

		nw, err := fp.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}

func (s *ServerGRPC) UploadFile(stream proto.RkUploaderService_UploadFileServer) (err error) {
	firstChunk := true
	var fp *os.File

	var files = make([]*proto.File, 0)

	var fileData *proto.UploadRequestType

	var filename string
	for {

		fileData, err = stream.Recv() //ignoring the data  TO-Do save files received
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}

		if firstChunk { //first chunk contains file name

			if fileData.Filename != "" { //create file

				fname := filepath.Base(fileData.Filename)

				//check if file ends with jpg or png or jpeg
				containsJpg := strings.Contains(fname, "jpg")
				containsGif := strings.Contains(fname, "gif")
				containsPng := strings.Contains(fname, "png")

				if !containsPng && !containsJpg && !containsGif {
					_ = stream.SendAndClose(&proto.UploadResponseType{
						Message: "file extension must be jpg, gif or png",
						Code:    proto.UploadStatusCode_Invalid,
					})
					return nil //errors.New("file extension must be jpg, jpeg or png")
				}

				//check if file exists
				if _, err := getFileByFileName(fname); err == nil {
					_ = stream.SendAndClose(&proto.UploadResponseType{
						Message: "File (" + fname + ") already exist",
						Code:    proto.UploadStatusCode_Exist,
					})
					return nil
				}

				fp, err = os.Create(path.Join(s.destDir, fname))
				if err != nil {
					s.logger.Error().Msg("Unable to create file  :" + fname)
					_ = stream.SendAndClose(&proto.UploadResponseType{
						Message: "Unable to create file :" + fname,
						Code:    proto.UploadStatusCode_Failed,
					})
					return
				}

				defer fp.Close()

			} else {
				s.logger.Error().Msg("FileName not provided in first chunk  :" + fileData.Filename)
				_ = stream.SendAndClose(&proto.UploadResponseType{
					Message: "FileName not provided in first chunk:" + fileData.Filename,
					Code:    proto.UploadStatusCode_Failed,
				})
				return

			}
			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {
			s.logger.Error().Msg("Unable to write chunk of filename :" + fileData.Filename + " " + err.Error())
			stream.SendAndClose(&proto.UploadResponseType{
				Message: "Unable to write chunk of filename :" + fileData.Filename,
				Code:    proto.UploadStatusCode_Failed,
			})
			return
		}

	}

	fn := filepath.Base(filename)
	files = append(files, &proto.File{FileName: fn, Uri: "data/files/" + fn})

	//s.logger.Info().Msg("upload received")
	err = stream.SendAndClose(&proto.UploadResponseType{
		Message: "Upload received with success",
		Code:    proto.UploadStatusCode_Ok,
		Files:   files,
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :" + filename + " in " + s.destDir)
	return
}
