package main

import (
	"fmt"
	"log"
	"os"
	"pick"

	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

func drawBox(x, y int) {
	if err := termbox.Clear(coldef, coldef); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	termbox.SetCell(x, y, '┏', coldef, coldef)
	termbox.SetCell(x+1, y, '┓', coldef, coldef)
	termbox.SetCell(x, y+1, '┗', coldef, coldef)
	termbox.SetCell(x+1, y+1, '┛', coldef, coldef)
	termbox.Flush() // Flushを呼ばないと描画されない
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s", err)
			os.Exit(1)
		}
	}()
	os.Exit(_main())
}

func _main() int {
	cli := pick.New()
	resp, err := cli.Run()

	if err != nil {
		return 1
	}

	fmt.Println(resp)

	return 0
}
