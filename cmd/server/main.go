package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"imageCache/delivery/server"
	"imageCache/pkg/env"
)

func main() {
	wait := make(chan int)
	_ = godotenv.Load()

	go func() {
		if er := server.StartGRPCServer(env.Get().GrpcAddr); er != nil {
			fmt.Println("=========>", er.Error())
			wait <- 1
			return
		}
	}()

	go func() {
		if er := server.StartRESTServer(env.Get().RestAddr); er != nil {
			fmt.Println("=========>", er.Error())
			wait <- 1
			return
		}
	}()

	//s := make(chan os.Signal, 1)
	//signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	//if n := <-s; n == syscall.SIGINT || n == syscall.SIGTERM || n == syscall.SIGHUP {
	//	wait <- 1
	//}

	<-wait
}
