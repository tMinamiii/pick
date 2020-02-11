package pick

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
)

type Pick struct{}

func New() *Pick {
	return &Pick{}
}

func (p *Pick) Run(term string) (*PocketGetResponse, error) {
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

	var key PocketAuthKey

	if json.Unmarshal(raw, &key) != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	request := NewPocketGetRequest(term, key)
	resp, err := request.Get()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
