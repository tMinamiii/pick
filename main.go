package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

// GetRequest is Pocket Retrieve API struct
type GetRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Search      string `json:"search"`
	Count       int    `json:"count"`
}

// NewGetRequest create new GetRequest data structure
func NewGetRequest(term string, key AuthKey) *GetRequest {
	return &GetRequest{
		ConsumerKey: key.ConsumerKey,
		AccessToken: key.AccessToken,
		Search:      term,
	}
}

func get(request *GetRequest) *PocketResponse {
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

	var presp PocketResponse
	if json.Unmarshal(byteArray, &presp) != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}

	return &presp
}

// OpenBrowser open url each platform default browser
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

// AuthKey is data structure for reading key.json
type AuthKey struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

func main() {
	// search("DynamoDB")
	raw, err := ioutil.ReadFile("./key.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var key AuthKey

	if json.Unmarshal(raw, &key) != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	request := NewGetRequest("DynamoDB", key)
	resp := get(request)
	fmt.Println(resp.String())
}
