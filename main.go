package main

import (
	"flag"
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

	// コマンドライン引数がある場合は、それを検索語とする
	flag.Parse()
	args := flag.Args()
	term := ""

	if len(args) > 0 {
		term = strings.Join(args, " ")
	}

	app.Action = func(context *cli.Context) error {
		return pocket.PickPocket(term)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
