package pick

import "bytes"

type PocketRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

// AuthKey is data structure for reading key.json
type PocketAuthKey struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}


// PocketResponse is represent Pocket Get API resposen
type PocketResponse struct {
	Status     int                       `json:"status"`
	List       map[string]*PocketArticle `json:"list"`
	Error      string                    `json:"error"`
	SearchMeta PocketSearchMeta          `json:"search_meta"`
	Since      int                       `json:"since"`
}

func (p *PocketResponse) String() string {
	var out bytes.Buffer

	for _, v := range p.List {
		out.WriteString(v.String())
	}

	return out.String()
}

// PocketSearchMeta represents
type PocketSearchMeta struct {
	TotalResultCount int    `json:"total_result_count"`
	Count            int    `json:"count"`
	Offset           int    `json:"offset"`
	HasMore          bool   `json:"has_more"`
	SearchType       string `json:"search_type"`
}

// PocketArticle represents Pocket Get API resposen
type PocketArticle struct {
	// A unique identifier matching the saved item.
	// This id must be used to perform any actions through the v3/modify endpoint.
	ItemID string `json:"item_id"`

	// A unique identifier similar to the item_id but is unique to the actual url of the saved item.
	// The resolved_id identifies unique urls. For example a direct link to a New York Times
	// article and a link that redirects (ex a shortened bit.ly url)
	// to the same article will share the same resolved_id.
	// If this value is 0, it means that Pocket has not processed the item.
	// Normally this happens within seconds but is possible you may request the item before it has been resolved.
	ResolveID string `json:"resolved_id"`

	// The actual url that was saved with the item. This url should be used if the user wants to view the item.
	GivenURL string `json:"given_url"`

	// The final url of the item. For example if the item was a shortened bit.ly link,
	// this will be the actual article the url linked to.
	ResolvedURL string `json:"resolved_url"`

	// The title that was saved along with the item.
	GivenTitle string `json:"given_title"`

	// The title that Pocket found for the item when it was parsed
	ResolvedTitle string `json:"resolved_title"`

	// 0 or 1 - 1 If the item is favorited
	Favorite string `json:"favorite"`

	// 0, 1, 2 - 1 if the item is archived - 2 if the item should be deleted
	Status string `json:"status"`

	// The first few lines of the item articles only)(
	Excerpt string `json:"excerpt"`

	// 0 or 1 - 1 if the item is an article
	IsArticle string `json:"is_article"`

	// 0, 1, or 2 - 1 if the item has images in it - 2 if the item is an image
	HasImage string `json:"has_image"`

	// 0, 1, or 2 - 1 if the item has videos in it - 2 if the item is a video
	HasVideo string `json:"has_video"`

	// How many words are in the article
	WordCount string `json:"word_count"`

	// A JSON object of the user tags associated with the item
	Tags string `json:"tags"`

	// A JSON object listing all of the authors associated with the item
	Authors string `json:"authors"`

	// A JSON object listing all of the images associated with the item videos)
	Images string `json:"images"`
}

func (p *PocketArticle) String() string {
	var out bytes.Buffer

	out.WriteString(" * " + p.ResolvedTitle)
	out.WriteString("(" + p.ResolvedURL + ")\n")

	return out.String()
}
