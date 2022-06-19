package cli

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "gopkg.in/cheggaaa/pb.v1"
	proto "imageCache/grpc/gen/proto/imageCache/v1"
)

const chunkSize = 64 * 1024 // 64 KiB

type uploader struct {
	dir         string
	client      proto.RkUploaderServiceClient
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	pool        *pb.Pool
	DoneRequest chan string
	FailRequest chan string
	addr        string
}

//NewUploader creates a object of type uploader and creates fixed worker goroutines/threads
func NewUploader(ctx context.Context, client proto.RkUploaderServiceClient, dir, addr string) *uploader {
	d := &uploader{
		ctx:         ctx,
		client:      client,
		dir:         dir,
		requests:    make(chan string),
		DoneRequest: make(chan string),
		FailRequest: make(chan string),
		addr:        addr,
	}
	for i := 0; i < 5; i++ {
		d.wg.Add(1)
		go d.worker(i + 1)
	}
	d.pool, _ = pb.StartPool()
	return d
}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}

func (d *uploader) worker(workerID int) {
	defer d.wg.Done()
	var (
		buf        []byte
		firstChunk bool
	)
	for request := range d.requests {

		//open
		//.Println("Processsing " + request)
		file, errOpen := os.Open(request)
		if errOpen != nil {
			errOpen = errors.Wrapf(errOpen,
				"failed to open file %s",
				request)
			return
		}

		defer file.Close()

		//start uploader
		streamUploader, err := d.client.UploadFile(d.ctx)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create upload stream for file %s",
				request)
			return
		}
		defer streamUploader.CloseSend()
		stat, errstat := file.Stat()
		if errstat != nil {
			err = errors.Wrapf(err,
				"Unable to get file size  %s",
				request)
			return
		}

		//start progress bar
		bar := pb.New64(stat.Size()).Postfix(" " + filepath.Base(request))
		bar.Units = pb.U_BYTES
		d.pool.Add(bar)

		//create a buffer of chunkSize to be streamed
		buf = make([]byte, chunkSize)
		firstChunk = true
		for {
			n, errRead := file.Read(buf)
			if errRead != nil {
				if errRead == io.EOF {
					errRead = nil
					break
				}

				errRead = errors.Wrapf(errRead,
					"errored while copying from file to buf")
				return
			}
			if firstChunk {
				err = streamUploader.Send(&proto.UploadRequestType{
					Content:  buf[:n],
					Filename: request,
				})
				firstChunk = false
			} else {
				err = streamUploader.Send(&proto.UploadRequestType{
					Content: buf[:n],
				})
			}
			if err != nil {

				bar.FinishPrint("failed to send chunk via stream file : " + request)
				break
				//bar.Reset(0)
				//return
			}

			bar.Add64(int64(n))
		}
		status, err := streamUploader.CloseAndRecv()

		if err != nil { //retry needed

			fmt.Println("failed to receive upstream status response")
			bar.FinishPrint("Error uploading file : " + request + " Error " + err.Error())
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		if status.Code == proto.UploadStatusCode_Exist {
			bar.FinishPrint("Error uploading file : " + request + " | " + status.Message)
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		if status.Code == proto.UploadStatusCode_Invalid {
			bar.FinishPrint("Error uploading file : " + request + " | " + status.Message)
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		if status.Code != proto.UploadStatusCode_Ok { //retry needed
			bar.FinishPrint("Error uploading file : " + request + " | " + status.Message)
			bar.Reset(0)
			d.FailRequest <- request
			return
		}

		fmt.Println("\nfile uploaded at: ")
		for _, file := range status.Files {
			fmt.Printf("{ FileName: \"%s\", Url: \"%s\" }\n", file.GetFileName(), file.GetUri())
		}

		//fmt.Println("writing done for : " + request + " by " + strconv.Itoa(workerID))
		d.DoneRequest <- request
		bar.Finish()
	}

}

func (d *uploader) Do(filepath string) {
	d.requests <- filepath
}

//UploadFiles takes in client grpcCLient object and  an optional list of file path or dir name as input.
//It sends the files  using the grpcClient object to the server in parallel
//returns error if file transfer is not successful
func UploadFiles(ctx context.Context, client proto.RkUploaderServiceClient, filepathlist []string, dir, addr string) error {

	d := NewUploader(ctx, client, dir, addr)
	var errorUploadBulk error

	if dir != "" {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		defer d.Stop()

		go func() {
			for _, file := range files {

				if !file.IsDir() {

					d.Do(dir + "/" + file.Name())

				}
			}
		}()

		for _, file := range files {
			if !file.IsDir() {
				select {

				case <-d.DoneRequest:

					//fmt.Println("sucessfully sent :" + req)

				case req := <-d.FailRequest:

					fmt.Println("failed to  send " + req)
					errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)

				}
			}
		}
		fmt.Println("All done")
	} else {

		go func() {
			for _, file := range filepathlist {
				d.Do(file)
			}
		}()

		defer d.Stop()

		for i := 0; i < len(filepathlist); i++ {
			select {

			case req := <-d.DoneRequest:
				fmt.Println("sucessfully sent " + req)
			case req := <-d.FailRequest:
				fmt.Println("failed to  send " + req)
				errorUploadBulk = errors.Wrapf(errorUploadBulk, " Failed to send %s", req)
			}
		}

	}

	return errorUploadBulk
}

func UploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "Uplooads files from server in parallel",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "cmd",
				Value: "command",
				Usage: "list",
			},
			cli.StringFlag{
				Name:  "a",
				Value: "localhost:port",
				Usage: "server address",
			},
			cli.StringFlag{
				Name:  "d",
				Value: "",
				Usage: "base directory",
			},
			cli.StringSliceFlag{
				Name:  "f",
				Usage: "files e.g -f sample1.png -f sample2.jpeg -f sample3.png",
				Value: nil,
			},
			cli.StringFlag{
				Name:  "tls-path",
				Value: "",
				Usage: "directory to the TLS server.crt file",
			},
		},
		Action: func(c *cli.Context) error {
			options := []grpc.DialOption{}
			if p := c.String("tls-path"); p != "" {
				creds, err := credentials.NewClientTLSFromFile(
					filepath.Join(p, "server.crt"),
					"")
				if err != nil {
					log.Println(err)
					return err
				}
				options = append(options, grpc.WithTransportCredentials(creds))
			} else {
				options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
			}
			addr := c.String("a")

			conn, err := grpc.Dial(addr, options...)
			if err != nil {
				log.Fatalf("cannot connect: %v", err)
			}
			defer conn.Close()

			files := c.StringSlice("f")

			return UploadFiles(context.Background(), proto.NewRkUploaderServiceClient(conn), files, c.String("d"), addr)
		},
	}
}
