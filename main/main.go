package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"pick/pocket"
)

func main() {
	// search("DynamoDB")
	raw, err := ioutil.ReadFile("./key.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var key pocket.PocketAuthKey

	if json.Unmarshal(raw, &key) != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	request := pocket.NewGetRequest("DynamoDB", key)
	resp := request.Get()
	fmt.Println(resp.String())
}
