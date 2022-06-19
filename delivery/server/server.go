package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"imageCache/data/files"
	proto "imageCache/grpc/gen/proto/imageCache/v1"
	"log"
	"net"
	"net/http"
	"os"
)

type Server interface {
	Listen() (err error)
	Close()
}
type ServerGRPC struct {
	logger  zerolog.Logger
	server  *grpc.Server
	Address string

	certificate string
	key         string
	destDir     string
	proto.UnimplementedRkUploaderServiceServer
}

type ServerGRPCConfig struct {
	Certificate string
	Key         string
	Address     string
	DestDir     string
}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	s.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "server").
		Logger()

	if cfg.Address == "" {
		err = errors.Errorf("Address must be specified")
		s.logger.Error().Msg("Address must be specified")
		return
	}

	s.Address = cfg.Address
	s.certificate = cfg.Certificate
	s.key = cfg.Key

	if _, err = os.Stat(cfg.DestDir); err != nil {
		s.logger.Error().Msg("Directory doesnt exist")
		return
	}

	s.destDir = cfg.DestDir

	return
}

func (s *ServerGRPC) Listen() (err error) {
	var (
		listener  net.Listener
		grpcOpts  = []grpc.ServerOption{}
		grpcCreds credentials.TransportCredentials
	)

	listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on  %s",
			s.Address)
		return
	}

	if s.certificate != "" && s.key != "" {
		grpcCreds, err = credentials.NewServerTLSFromFile(
			s.certificate, s.key)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create tls grpc server using cert %s and key %s",
				s.certificate, s.key)
			return
		}

		grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	}

	s.server = grpc.NewServer(grpcOpts...)
	proto.RegisterRkUploaderServiceServer(s.server, s)

	err = s.server.Serve(listener)
	if err != nil {
		err = errors.Wrapf(err, "errored listening for grpc connections")
		return
	}

	log.Println("server running @", s.Address)

	return
}

func (s *ServerGRPC) Close() {
	if s.server != nil {
		s.server.Stop()
	}

	log.Println("server @", s.Address, " shutdown")
	return
}

func StartGRPCServer(address string) error {

	grpcServer, err := NewServerGRPC(ServerGRPCConfig{
		Address: address,
		DestDir: files.GetLocation(), //c.String("d"),
	})

	if err != nil {
		fmt.Println("error is creating server =========> ", err.Error())
		return err
	}
	server := &grpcServer

	fmt.Println("=========> We are running GRPC @", address)

	err = server.Listen()
	if err != nil {
		fmt.Println("=========> ", err.Error())
	}

	defer server.Close()
	fmt.Println("=========> We are stopping")
	return nil
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func StartRESTServer(address string) error {

	engine := SetUpRouter()
	engine.GET("data/files/:fileName", DownloadFileHandler)
	engine.GET("list", ListFilesHandler)

	fmt.Println("======> we are running REST @", address)

	return http.ListenAndServe(address, engine)
}
