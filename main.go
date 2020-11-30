package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	flag.Parse()
	args := flag.Args()
	word := ""

	if len(args) > 0 {
		word = strings.Join(args, " ")
	}

	fmt.Println(word)

	app.Action = func(context *cli.Context) error {
		return pocket.PickPocket(word)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
