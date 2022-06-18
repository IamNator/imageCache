package main

import (
	"log"
	"os"
	"time"

	"imageCache/delivery/server"
)

func main() {
	//app := cli.NewApp()
	//app.Name = "rk_mft"
	//app.Usage = "RK multi File Transferer Server"
	//app.Version = "0.0.1"
	myApp := server.StartServerCommand()

	//app.Setup()

	if err := myApp.Run(os.Args); err != nil {
		log.Println("Stop.", err.Error())
	}

	log.Println("========> Ran this")
	time.Sleep(time.Hour)
}
