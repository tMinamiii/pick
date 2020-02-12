// pocket package
package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"

	"github.com/manifoldco/promptui"
)

type Request struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

type AuthKey struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

type SearchMeta struct {
	TotalResultCount int    `json:"total_result_count"`
	Count            int    `json:"count"`
	Offset           int    `json:"offset"`
	HasMore          bool   `json:"has_more"`
	SearchType       string `json:"search_type"`
}

type Article struct {
	ItemID        string `json:"item_id"`
	ResolveID     string `json:"resolved_id"`
	GivenURL      string `json:"given_url"`
	ResolvedURL   string `json:"resolved_url"`
	GivenTitle    string `json:"given_title"`
	ResolvedTitle string `json:"resolved_title"`
	Favorite      string `json:"favorite"`
	Status        string `json:"status"`
	Excerpt       string `json:"excerpt"`
	IsArticle     string `json:"is_article"`
	HasImage      string `json:"has_image"`
	HasVideo      string `json:"has_video"`
	WordCount     string `json:"word_count"`
	Tags          string `json:"tags"`
	Authors       string `json:"authors"`
	Images        string `json:"images"`
}

func (p *Article) String() string {
	// var out bytes.Buffer
	// out.WriteString(p.ResolvedURL)
	// return out.String()
	return p.ResolvedURL
}

type GetResponse struct {
	Status     int                 `json:"status"`
	List       map[string]*Article `json:"list"`
	Error      string              `json:"error"`
	SearchMeta SearchMeta          `json:"search_meta"`
	Since      int                 `json:"since"`
}

func (p *GetResponse) String() string {
	var out bytes.Buffer

	for _, v := range p.List {
		out.WriteString(v.String())
	}

	return out.String()
}

type GetRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Search      string `json:"search"`
	Count       int    `json:"count"`
}

func NewPocketGetRequest(term string, key AuthKey) *GetRequest {
	return &GetRequest{
		ConsumerKey: key.ConsumerKey,
		AccessToken: key.AccessToken,
		Search:      term,
		Count:       100,
	}
}

func (request *GetRequest) Get() (*GetResponse, error) {
	url := "https://getpocket.com/v3/get"

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to marshal struct object. %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var presp GetResponse

	err = json.Unmarshal(byteArray, &presp)
	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		return nil, err
	}

	return &presp, nil
}

func PickPocket() (err error) {
	prompt := promptui.Prompt{
		Label: "Search",
	}

	term, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	usr, err := user.Current()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	raw, err := ioutil.ReadFile(usr.HomeDir + "/.config/pick/key.json")

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var key AuthKey

	if err = json.Unmarshal(raw, &key); err != nil {
		log.Fatal(err.Error())
		return
	}

	request := NewPocketGetRequest(term, key)
	resp, err := request.Get()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	items := make([]*Article, 0, len(resp.List))
	for _, val := range resp.List {
		items = append(items, val)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "* {{ .ResolvedTitle | red }}",
		Inactive: " {{ .ResolvedTitle | cyan }}",
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
		return
	}

	OpenBrowser(url)

	return nil
}
