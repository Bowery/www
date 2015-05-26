package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

const usage = `Usage: www <provider> [options]
www reads from standard input and pipes the given
input to a provider.

Providers: slack, gist.

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
	// Verify there is valid input from Stdin.
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		log.Fatal(err)
		os.Exit(1)
	}

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
			log.Fatal(err)
		}
	}

	err = db.Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Validate and set the provider.
	if len(os.Args) < 1 {
		log.Fatal("Provider required.")
	}

	provider := ps[os.Args[1]]
	if provider == nil {
		log.Fatal("Invalid provider.")
	}

	// Fill the config with an empty entry if nil.
	_, ok := config[os.Args[1]]
	if !ok {
		config[os.Args[1]] = map[string]string{}
	}

	// Read in Stdin.
	var content bytes.Buffer
	_, err = io.Copy(&content, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Initialize provider. Passes flags and config
	// for that specific provider.
	err = provider.Init(os.Args[2:], config[os.Args[1]])
	if err != nil {
		log.Fatal(err)
	}

	// Execute `Send` method of provider.
	err = provider.Send(content)
	if err != nil {
		panic(err)
	}

	// Save any configuration changes made by provider.
	err = db.Save(config)
	if err != nil {
		log.Fatal(err)
	}
}
