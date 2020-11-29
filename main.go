package main

import (
	"log"
	"os"

	"github.com/tMinamiii/pick/pocket"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pick"
	app.Usage = "This app search sites in your Pocket and open browser"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "keygen",
			Usage: "Generate key.json in ~/.config/pick",
			Action: func(c *cli.Context) error {
				pocket.RunKeyGen()

				return nil
			},
		},
	}

	app.Action = func(context *cli.Context) error {
		return pocket.PickPocket()
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
