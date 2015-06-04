package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Bowery/conf"
	"github.com/Bowery/www/providers"
)

var (
	configPath string
	err        error
	db         *conf.JSON
	ps         map[string]providers.Provider
)

const usage = `Usage:
  $ somecmd | www <provider> [options]

www reads from standard input and pipes the given
input to a provider.

Providers: slack, gist, gmail, s3.

To setup a provider
  $ www setup <provider>

For information on how to use a provider
  $ www <provider> --help
`

func init() {
	homeVar := "HOME"
	if runtime.GOOS == "windows" {
		homeVar = "USERPROFILE"
	}
	configPath = filepath.Join(os.Getenv(homeVar), ".wwwconf")
	ps = providers.Providers
}

func main() {
	// Create new local config and attempt to read in configuration
	// file specified by configPath. If the file does not exist,
	// create it.
	config := map[string]map[string]string{}
	db, err = conf.NewJSON(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			dat, _ := json.Marshal(config)
			ioutil.WriteFile(configPath, dat, os.ModePerm)
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	err = db.Load(&config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Validate and set the provider.
	if len(os.Args) <= 1 {
		fmt.Println(usage)
		os.Exit(2)
	}

	var provider providers.Provider
	if os.Args[1] == "setup" {
		if len(os.Args) > 2 {
			_, ok := config[os.Args[2]]
			if !ok {
				config[os.Args[2]] = map[string]string{}
			}

			provider = ps[os.Args[2]]
			provider.Setup(config[os.Args[2]])
			err = db.Save(config)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}

			os.Exit(0)
		}
	}

	provider = ps[os.Args[1]]
	if provider == nil {
		fmt.Fprintln(os.Stderr, "Invalid provider.")
		os.Exit(1)
	}

	// Fill the config with an empty entry if nil.
	_, ok := config[os.Args[1]]
	if !ok {
		config[os.Args[1]] = map[string]string{}
	}

	// Initialize provider. Passes flags and config
	// for that specific provider.
	err = provider.Init(os.Args[2:], config[os.Args[1]])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Verify there is valid input from Stdin.
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// If there is no valid input print help.
	if stat.Mode()&os.ModeCharDevice != 0 {
		os.Exit(1)
	}

	// Read in Stdin.
	var content bytes.Buffer
	_, err = io.Copy(&content, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Execute `Send` method of provider.
	err = provider.Send(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Save any configuration changes made by provider.
	err = db.Save(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
