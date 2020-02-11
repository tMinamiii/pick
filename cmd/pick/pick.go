package main

import (
	"flag"
	"fmt"
	"os"
	"pick"

	"github.com/manifoldco/promptui"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 2 {
		if args[0] == "keygen" {
			consumerKey := args[1]
			code := pick.RequestCode(consumerKey)
			fmt.Println(code)
			// open browser for auth
			go pick.AuthAndRedirect(code)
			// launch http server for detecting auth finished
			pick.LaunchHTTPServer(consumerKey, code)
		}

		os.Exit(0)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s", err)
			os.Exit(1)
		}
	}()
	os.Exit(_main())
}

func _main() int {
	prompt := promptui.Prompt{
		Label: "Search",
	}

	term, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	cli := pick.New()
	resp, err := cli.Run(term)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	items := make([]*pick.PocketArticle, 0, len(resp.List))
	for _, val := range resp.List {
		items = append(items, val)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "> {{ .ResolvedTitle | cyan }}",
		Inactive: "{{ .ResolvedTitle | cyan }}",
		// Active:   "> {{ .ResolvedTitle | cyan }} ({{ .ResolvedURL | red }})",
		// Inactive: "{{ .ResolvedTitle | cyan }} ({{ .ResolvedURL | red }})",
		// Selected: "> {{ .ResolvedTitle | red | cyan }}",
	}
	selectPrompt := promptui.Select{
		Label:     "Select Site",
		Size:      30,
		Items:     items,
		Templates: templates,
	}
	_, url, err := selectPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	pick.OpenBrowser(url)

	return 0
}
