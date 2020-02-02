package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pick"
	"strings"
	"time"
)

const (
	RedirectURL = "localhost:18123"
)

func Post(url string, payload []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		byteArray, _ := ioutil.ReadAll(resp.Body)

		return byteArray
	}

	return []byte{}
}

func requestCode(consumerKey string) string {
	request := struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}{
		consumerKey,
		"http://" + RedirectURL,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to marshal struct object. %v\n", err)
		os.Exit(1)
	}

	url := "https://getpocket.com/v3/oauth/request"
	respBody := Post(url, payload)

	return strings.Split(string(respBody), "=")[1]
}

func authAndRedirect(code string) {
	authorizeRequestURL := "https://getpocket.com/auth/authorize?request_token=" + code + "&redirect_uri=http://" + RedirectURL

	time.Sleep(1 * time.Second)
	pick.OpenBrowser(authorizeRequestURL)
}

func emitAccessToken(consumerKey string, code string) []byte {
	authorize := struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}{
		consumerKey,
		code,
	}

	payload, err := json.Marshal(authorize)
	if err != nil {
		fmt.Printf("Failed to marshal struct object. %v\n", err)
		os.Exit(1)
	}

	authURL := "https://getpocket.com/v3/oauth/authorize"
	body := Post(authURL, payload)

	fmt.Println(string(body))

	if len(body) > 0 {
		accessTokenParameter := strings.Split(string(body), "&")[0]
		accessToken := strings.Split(accessTokenParameter, "=")[1]

		result, err := json.Marshal(&pick.PocketAuthKey{
			ConsumerKey: consumerKey,
			AccessToken: accessToken,
		})

		if err != nil {
			os.Exit(1)
		}

		return result
	}

	return []byte{}
}

type redirectHandler struct {
	ConsumerKey string
	Code        string
}

func (rh redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := emitAccessToken(rh.ConsumerKey, rh.Code)

	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "<script type=\"text/javascript\">")
	fmt.Fprintln(w, "fetch('http://localhost:18123/exit');")
	fmt.Fprintln(w, "</script>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>authorized</body>")
	fmt.Fprintln(w, "</html>")

	if len(result) > 0 {
		file, err := os.Create("./key.json")

		if err != nil {
			fmt.Println("Failed to create key.json")
			os.Exit(1)
		}

		defer file.Close()

		_, err = file.Write(result)

		if err != nil {
			fmt.Println("Failed to write key.json")
			os.Exit(1)
		}
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()

	consumerKey := args[0]
	code := requestCode(consumerKey)

	go authAndRedirect(code)

	http.Handle("/", &redirectHandler{consumerKey, code})
	http.HandleFunc("/exit", exitHandler)

	if err := http.ListenAndServe(RedirectURL, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
