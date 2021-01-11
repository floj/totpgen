package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

type configEntry struct {
	Name   string
	Secret string
}

type config []configEntry

func (c config) find(name string) *configEntry {
	for _, e := range c {
		if e.Name == name {
			return &e
		}
	}
	return nil
}

const usage = `
Usage:
 Call with secret name
   %s <secret-name>
  or symlink to 'totpgen-secret-name'
`

func main() {
	token, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(token)

}

func run() (string, error) {
	confName := getConfName(os.Args)
	if confName == "" {
		return "", fmt.Errorf(usage, os.Args[0])
	}
	conf, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("Could not open config: %w", err)
	}

	entry := conf.find(confName)
	if entry == nil {
		return "", fmt.Errorf("No config entry for '%s' found", confName)
	}

	token, err := totp.GenerateCode(entry.Secret, time.Now())
	return token, err
}

func loadConfig() (*config, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	confFile := filepath.Join(confDir, "totpgen", "config.json")
	f, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &config{}
	err = json.NewDecoder(f).Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getConfName(args []string) string {
	cmd := os.Args[0]
	if idx := strings.Index(cmd, "-"); idx >= 0 {
		return cmd[idx+1:]
	}
	if len(args) == 2 {
		return args[1]
	}
	return ""
}
