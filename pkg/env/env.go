package env

import (
	"github.com/joho/godotenv"
	"os"
)

type (
	Env struct {
		GrpcAddr string
		RestAddr string
	}
)

func Get() Env {
	if &env == nil {
		InitEnv()
	}

	if env.RestAddr == "" {
		InitEnv()
	}

	return env
}

var env Env

func InitEnv() {
	godotenv.Load()

	env = Env{
		GrpcAddr: os.Getenv("GRPC_ADDR"),
		RestAddr: os.Getenv("REST_ADDR"),
	}
}
