package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	pickpocket "pick/ppt"
)

// GetRequest is Pocket Retrieve API struct
type GetRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Search      string `json:"search"`
	Count       int    `json:"count"`
}

// NewGetRequest create new GetRequest data structure
func NewGetRequest(term string, key pick.PocketAuthKey) *GetRequest {
	return &GetRequest{
		ConsumerKey: key.ConsumerKey,
		AccessToken: key.AccessToken,
		Search:      term,
	}
}

func (request *GetRequest) get() *pickpocket.PocketResponse {
	url := "https://getpocket.com/v3/get"

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to marshal struct object. %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var presp pick.PocketResponse
	if json.Unmarshal(byteArray, &presp) != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}

	return &presp
}

func main() {
	// search("DynamoDB")
	raw, err := ioutil.ReadFile("./key.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var key pick.PocketAuthKey

	if json.Unmarshal(raw, &key) != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	request := NewGetRequest("DynamoDB", key)
	resp := request.get()
	fmt.Println(resp.String())
}
