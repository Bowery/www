package providers

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/smtp"

	"github.com/Bowery/prompt"
)

func init() {
	Providers["gmail"] = &GMail{}
}

// GMail represents a GMail provider.
type GMail struct {
	flag     *flag.FlagSet
	subject  string
	to       string
	user     string
	password string
}

// Init reads in flags, setting the recipient and authentication.
func (g *GMail) Init(args []string, config map[string]string) error {
	g.flag = flag.NewFlagSet("gmail", flag.ExitOnError)
	g.flag.StringVar(&g.subject, "subject", "stdout via www", "Subject of message.")
	g.flag.StringVar(&g.to, "to", "", "Recipient of message.")
	err := g.flag.Parse(args)
	if err != nil {
		return err
	}

	// Recipient required.
	if g.to == "" {
		return errors.New("recipient required")
	}

	// Prompt user for email address and password.
	if config["user"] != "" {
		g.user = config["user"]
	} else {
		g.user, err = prompt.Basic("Email: ", true)
		if err != nil {
			return err
		}

		config["user"] = g.user
	}

	if config["password"] != "" {
		g.password = config["password"]
	} else {
		g.password, err = prompt.Password("Password: ")
		if err != nil {
			return err
		}

		config["password"] = g.password
	}

	return nil
}

// Send sends the email with the stdout as the body.
func (g *GMail) Send(content bytes.Buffer) error {
	body := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", g.to, g.subject, content.String())
	auth := smtp.PlainAuth("", g.user, g.password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, g.user, []string{g.to}, []byte(body))
	if err != nil {
		return err
	}

	return nil
}
