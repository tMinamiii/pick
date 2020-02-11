package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/manifoldco/promptui"
	"github.com/tMinamiii/pick/pocket"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 2 {
		if args[0] == "keygen" {
			consumerKey := args[1]
			code := pocket.RequestCode(consumerKey)
			fmt.Println(code)
			// open browser for auth
			go pocket.AuthAndRedirect(code)
			// launch http server for detecting auth finished
			pocket.LaunchHTTPServer(consumerKey, code)
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

	resp, err := run(term)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 1
	}

	items := make([]*pocket.Article, 0, len(resp.List))
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

	pocket.OpenBrowser(url)

	return 0
}

func run(term string) (*pocket.GetResponse, error) {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	raw, err := ioutil.ReadFile(usr.HomeDir + "/.config/pick/key.json")

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	var key pocket.AuthKey

	if json.Unmarshal(raw, &key) != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	request := pocket.NewPocketGetRequest(term, key)
	resp, err := request.Get()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
