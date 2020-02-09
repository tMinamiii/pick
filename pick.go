package pick

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Pick struct {
	Argv   []string
	Stderr io.Writer
	Stdin  io.Reader
	Stdout io.Writer
}

func newIDGen() *idgen {
	return &idgen{ch: make(chan uint64)}
}

func (ig *idgen) Run(ctx context.Context) {
	var i uint64
	for ; ; i++ {
		select {
		case <-ctx.Done():
			return
		case ig.ch <- i:
		}

		if i >= uint64(1<<63)-1 {
			// If this happens, it's a disaster, but what can we do...
			i = 0
		}
	}
}

func (ig *idgen) Next() uint64 {
	return <-ig.ch
}

func New() *Pick {
	return &Pick{
		Argv:   os.Args,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
}

func (p *Pick) Run() (string, error) {
	term := p.Argv[0]

	raw, err := ioutil.ReadFile("./key.json")
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}

	var key PocketAuthKey

	if json.Unmarshal(raw, &key) != nil {
		log.Fatal(err.Error())
		return "", err
	}

	request := NewPocketGetRequest(term, key)
	resp, err := request.Get()
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
