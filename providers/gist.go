package providers

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func init() {
	Providers["gist"] = &Gist{}
}

var (
	gistsAPIEndpoint = "https://api.github.com/gists"
)

// Gist represents a GitHub Gist provider. By default
// anonymous Gist's are generated.
type Gist struct {
	flag *flag.FlagSet
	name string
}

type gistReqBody struct {
	Public bool                         `json:"public"`
	Files  map[string]map[string]string `json:"files"`
}

type gistResBody struct {
	HTMLURL string `json:"html_url"`
}

// Setup does nothing for now.
func (g *Gist) Setup(config map[string]string) error {
	return nil
}

// Init reads in flags and sets the file name if specified.
func (g *Gist) Init(args []string, config map[string]string) error {
	g.flag = flag.NewFlagSet("gist", flag.ExitOnError)
	g.flag.StringVar(&g.name, "name", "", "Name of file.")
	err := g.flag.Parse(args)
	if err != nil {
		return err
	}

	if g.name == "" {
		g.name = fmt.Sprintf("stdgist-%d", time.Now().Unix())
	}

	return nil
}

// Send posts the content to Gist's api.
func (g *Gist) Send(content bytes.Buffer) error {
	body := gistReqBody{
		Public: true,
		Files: map[string]map[string]string{
			g.name: map[string]string{
				"content": content.String(),
			},
		},
	}

	content.Reset()
	encoder := json.NewEncoder(&content)
	err := encoder.Encode(body)
	if err != nil {
		return err
	}

	res, err := http.Post(gistsAPIEndpoint, "application/json", &content)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var resBody gistResBody
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resBody)
	if err != nil {
		return err
	}

	if resBody.HTMLURL != "" {
		fmt.Println(resBody.HTMLURL)
	}

	return nil
}
