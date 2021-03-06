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

	"github.com/manifoldco/promptui"
)

const (
	redirectURL               = "localhost:18123"
	oauthRequestURL           = "https://getpocket.com/v3/oauth/request"
	oauthAuthorizeURL         = "https://getpocket.com/v3/oauth/authorize"
	authorizeRequestURLFormat = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=http://localhost:18123"
)

func RunKeyGen() {
	prompt := promptui.Prompt{
		Label: "Input Consumer Key",
	}

	consumerKey, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	code := RequestCode(consumerKey)
	// open browser for auth
	go AuthAndRedirect(code)
	// launch http server for detecting auth finished
	LaunchHTTPServer(consumerKey, code)
}

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

	respBody := Post(oauthRequestURL, payload)

	return strings.Split(string(respBody), "=")[1]
}

func AuthAndRedirect(code string) {
	authorizeRequestURL := fmt.Sprintf(authorizeRequestURLFormat, code)
	// "https://getpocket.com/auth/authorize?request_token=" + code + "&redirect_uri=http://" + redirectURL

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

	body := Post(oauthAuthorizeURL, payload)

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

	html := `<html>
<head>
<script type="text/javascript">
fetch('http://localhost:18123/exit');
</script>
</head>
<body>authorized</body>
</html>
`
	fmt.Fprintln(w, html)

	if len(result) > 0 {
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err.Error())
			return
		}

		configDir := usr.HomeDir + "/.config"
		pickDir := configDir + "/pick"

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			os.MkdirAll(pickDir, 0755)
		} else if _, err := os.Stat(pickDir); os.IsNotExist(err) {
			os.MkdirAll(pickDir, 0755)
		}

		file, err := os.Create(pickDir + "/key.json")

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

func LaunchHTTPServer(consumerKey string, code string) {
	http.Handle("/", &redirectHandler{ConsumerKey: consumerKey, Code: code})
	http.HandleFunc("/exit", exitHandler)

	if err := http.ListenAndServe(redirectURL, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
