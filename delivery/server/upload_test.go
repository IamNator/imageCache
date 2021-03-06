package server_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"imageCache/data/files"
	"imageCache/delivery/server"
	proto "imageCache/grpc/gen/proto/imageCache/v1"
	"imageCache/pkg/env"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	env.InitEnv()

	grpcServer, err := server.NewServerGRPC(server.ServerGRPCConfig{
		Address: env.Get().GrpcAddr,
		DestDir: files.GetLocation(), //c.String("d"),
	})
	if err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}

	proto.RegisterRkUploaderServiceServer(s, &grpcServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	fmt.Println("===> big man is running here!!")
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestUploadFile(t *testing.T) {

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewRkUploaderServiceClient(conn)
	uploadClient, err := client.UploadFile(ctx)
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}

	fname := "HiFile.png"
	uploadClient.Send(&proto.UploadRequestType{
		Content:  []byte("hi how are you"),
		Filename: fname,
	})

	resp, err := uploadClient.CloseAndRecv()
	if err != nil {
		t.Fatalf("close failed failed: %v", err)
	}

	if er := os.Remove(files.GetLocation() + fname); er != nil {
		t.Fatalf(er.Error())
	}

	log.Printf("Response: %+v", resp)
}
