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

	godotenv.Load("../../.env", ".") //for during tests and all

	env = Env{
		GrpcAddr: mustGet("GRPC_ADDR"),
		RestAddr: mustGet("REST_ADDR"),
	}
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " is required in env")
	}
	return v
}
