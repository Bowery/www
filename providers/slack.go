package providers

import (
	"bytes"
	"errors"
	"flag"

	"github.com/Bowery/prompt"
	"github.com/Bowery/slack"
)

func init() {
	Providers["slack"] = &Slack{}
}

var (
	errSlackChannelRequired = errors.New("Channel required.")
)

// Slack represents a Slack provider.
type Slack struct {
	flag    *flag.FlagSet
	client  *slack.Client
	channel string
	token   string
	user    string
}

// Setup asks the user for their slack token.
func (s *Slack) Setup(config map[string]string) error {
	var err error
	config["token"], err = prompt.Basic("Token: ", true)
	if err != nil {
		return err
	}

	return nil
}

// Init reads in flags, setting the channel and token.
func (s *Slack) Init(args []string, config map[string]string) error {
	s.flag = flag.NewFlagSet("slack", flag.ExitOnError)
	s.flag.StringVar(&s.channel, "channel", "", "The slack #channel to post to.")
	s.flag.StringVar(&s.token, "token", "", "Authorization token.")
	s.flag.StringVar(&s.user, "user", "www", "User name.")
	err := s.flag.Parse(args)
	if err != nil {
		return err
	}

	// If no token passed as flag, check the config. If no
	// token has been sent, prompt the user for a token and
	// set on config variable. This will allow the token
	// to persist to next usage.
	if s.token == "" {
		if config["token"] != "" {
			s.token = config["token"]
		}
	}

	if s.token == "" {
		return errors.New("must set up.")
	}

	// Channel required.
	if s.channel == "" {
		return errSlackChannelRequired
	}

	// Create slack client.
	s.client = slack.NewClient(s.token)
	return nil
}

// Send posts the message to Slack.
func (s *Slack) Send(content bytes.Buffer) error {
	if s.channel[0] != '#' {
		s.channel = "#" + s.channel
	}

	return s.client.SendMessage(s.channel, content.String(), s.user)
}
