package main

import (
	"fmt"
	"os"
	"pick"
)

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
