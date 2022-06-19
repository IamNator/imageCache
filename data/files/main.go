package main

import (
	"os"

	"github.com/urfave/cli"
	app "imageCache/delivery/cli"
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "rk_mft"
	myApp.Usage = "RK multi File Transferer"
	myApp.Version = "0.0.1"
	myApp.Commands = []cli.Command{
		app.UploadCommand(),
	}
	myApp.Run(os.Args)
}
