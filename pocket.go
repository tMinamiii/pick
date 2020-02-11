package pick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

type PocketAuthKey struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

type PocketSearchMeta struct {
	TotalResultCount int    `json:"total_result_count"`
	Count            int    `json:"count"`
	Offset           int    `json:"offset"`
	HasMore          bool   `json:"has_more"`
	SearchType       string `json:"search_type"`
}

type PocketArticle struct {
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

func (p *PocketArticle) String() string {
	// var out bytes.Buffer
	// out.WriteString(p.ResolvedURL)
	// return out.String()
	return p.ResolvedURL
}

type PocketGetResponse struct {
	Status     int                       `json:"status"`
	List       map[string]*PocketArticle `json:"list"`
	Error      string                    `json:"error"`
	SearchMeta PocketSearchMeta          `json:"search_meta"`
	Since      int                       `json:"since"`
}

func (p *PocketGetResponse) String() string {
	var out bytes.Buffer

	for _, v := range p.List {
		out.WriteString(v.String())
	}

	return out.String()
}

type PocketGetRequest struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Search      string `json:"search"`
	Count       int    `json:"count"`
}

func NewPocketGetRequest(term string, key PocketAuthKey) *PocketGetRequest {
	return &PocketGetRequest{
		ConsumerKey: key.ConsumerKey,
		AccessToken: key.AccessToken,
		Search:      term,
		Count:       100,
	}
}

func (request *PocketGetRequest) Get() (*PocketGetResponse, error) {
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

	var presp PocketGetResponse

	err = json.Unmarshal(byteArray, &presp)
	if err != nil {
		fmt.Printf("Failed to create NewRequest. %v\n", err)
		return nil, err
	}

	return &presp, nil
}
