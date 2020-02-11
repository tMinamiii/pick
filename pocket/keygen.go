package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"
)

const (
	redirectURL = "localhost:18123"
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

func RequestCode(consumerKey string) string {
	request := struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}{
		consumerKey,
		"http://" + redirectURL,
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

func AuthAndRedirect(code string) {
	authorizeRequestURL := "https://getpocket.com/auth/authorize?request_token=" + code + "&redirect_uri=http://" + redirectURL

	time.Sleep(1 * time.Second)
	OpenBrowser(authorizeRequestURL)
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

		result, err := json.Marshal(&AuthKey{
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
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err.Error())
			return
		}

		file, err := os.Create(usr.HomeDir + "/.config/pick/key.json")
		defer file.Close()

		if err != nil {
			fmt.Println("Failed to create key.json")
			os.Exit(1)
		}

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

func LaunchHTTPServer(consumerKey string, code string) {
	http.Handle("/", &redirectHandler{ConsumerKey: consumerKey, Code: code})
	http.HandleFunc("/exit", exitHandler)

	if err := http.ListenAndServe(redirectURL, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
