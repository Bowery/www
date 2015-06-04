package providers

import (
	"bytes"
	"errors"
	"flag"
	"fmt"

	"github.com/Bowery/prompt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func init() {
	Providers["s3"] = &S3{}
}

// S3 represents an S3 provider.
type S3 struct {
	flag   *flag.FlagSet
	path   string
	bucket string
	access string
	secret string
	client *s3.S3
}

// Setup asks the user for their access and secret keys.
func (s *S3) Setup(config map[string]string) error {
	var err error
	config["access"], err = prompt.Basic("Access: ", true)
	if err != nil {
		return err
	}

	config["secret"], err = prompt.Basic("Secret: ", true)
	if err != nil {
		return err
	}

	return nil
}

// Init reads in flags, setting the access and secret keys.
func (s *S3) Init(args []string, config map[string]string) error {
	s.flag = flag.NewFlagSet("s3", flag.ExitOnError)
	s.flag.StringVar(&s.path, "path", "", "Name of file.")
	s.flag.StringVar(&s.bucket, "bucket", "", "Name of bucket.")
	s.flag.StringVar(&s.access, "access", "", "Access Key.")
	s.flag.StringVar(&s.secret, "secret", "", "Secret Key.")
	err := s.flag.Parse(args)
	if err != nil {
		return err
	}

	if config["access"] != "" {
		s.access = config["access"]
	}

	if config["secret"] != "" {
		s.secret = config["secret"]
	}

	if s.access == "" || s.secret == "" {
		return errors.New("must set up.")
	}

	if s.path == "" {
		return errors.New("path required.")
	}

	if s.bucket == "" {
		return errors.New("bucket required.")
	}

	auth, err := aws.GetAuth(s.access, s.secret)
	if err != nil {
		return err
	}

	s.client = s3.New(auth, aws.USEast)
	return nil
}

// Send puts the content to S3.
func (s *S3) Send(content bytes.Buffer) error {
	bucket := s.client.Bucket(s.bucket)
	err := bucket.PutReader(s.path, &content, int64(content.Len()), "text/plain", s3.ACL("public-read"))
	if err != nil {
		return err
	}

	fmt.Printf("http://%s.s3.amazonaws.com/%s", s.bucket, s.path)
	return nil
}
