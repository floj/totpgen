package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/pquerna/otp/totp"

	"github.com/urfave/cli/v2"
)

const ringName = "totpgen"

func main() {
	appBin := filepath.Base(os.Args[0])
	app := &cli.App{
		Name:  appBin,
		Usage: "prints a totp code for a configured secret",
		Commands: []*cli.Command{
			{
				Name:      "set",
				Usage:     "add or replace a configured secret",
				UsageText: "totpgen set <secret-name> <secret-token>",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() != 2 {
						cli.ShowSubcommandHelpAndExit(ctx, 1)
					}
					ring, err := openKeyring()
					if err != nil {
						return err
					}
					err = ring.Set(keyring.Item{
						Key:  ctx.Args().Get(0),
						Data: []byte(ctx.Args().Get(1)),
					})
					if err != nil {
						return fmt.Errorf("could not insert new secret into keyring: %w", err)
					}
					return nil
				},
			},
			{
				Name:      "list",
				Usage:     "show configured secrets",
				UsageText: "totpgen list",
				Action: func(ctx *cli.Context) error {
					ring, err := openKeyring()
					if err != nil {
						return err
					}
					kk, err := ring.Keys()
					if err != nil {
						return fmt.Errorf("could not list keys from keyring: %w", err)
					}
					sort.Strings(kk)
					for _, k := range kk {
						fmt.Println(k)
					}
					return nil
				},
			},
		},
		UsageText: "totpgen <secret-name>\ntotpgen-<secret-name> (eg. symlinked)\ntotpgen command [command options]",
		ArgsUsage: "<secret name>",
		Action: func(ctx *cli.Context) error {
			secName := ""
			if idx := strings.Index(appBin, "-"); idx >= 0 {
				secName = appBin[idx+1:]
			} else {
				secName = ctx.Args().First()
			}

			if secName == "" {
				cli.ShowAppHelpAndExit(ctx, 1)
			}

			ring, err := openKeyring()
			if err != nil {
				return err
			}
			e, err := ring.Get(secName)
			if err != nil {
				return fmt.Errorf("could not get secret for '%s': %w", secName, err)
			}

			token, err := totp.GenerateCode(string(e.Data), time.Now())
			if err != nil {
				return fmt.Errorf("could not generate totp code for '%s': %w", secName, err)
			}
			fmt.Print(token)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func openKeyring() (keyring.Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: ringName,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open keyring: %w", err)
	}
	return ring, nil
}

// func main() {
// 	addSecret := flag.Bool("set", false, "add/replace an OTP secret, use like -set <secret name> <secret token>")
// 	flag.Parse()

// 	token, err := run()
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// 	fmt.Print(token)
// }

// func run() (string, error) {
// 	ring, err := keyring.Open(keyring.Config{
// 		ServiceName: "github.com/floj/totpgen",
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("could not open keyring: %w", err)
// 	}

// 	confName := getConfName(os.Args)
// 	if confName == "" {
// 		return "", fmt.Errorf(usage, os.Args[0])
// 	}

// 	e, err := ring.Get(confName)
// 	if err != nil {
// 		return "", fmt.Errorf("could not get secret for '%s': %w", confName, err)
// 	}

// 	token, err := totp.GenerateCode(string(e.Data), time.Now())
// 	return token, err
// }

// func getConfName(args []string) string {
// 	cmd := os.Args[0]
// 	if idx := strings.Index(cmd, "-"); idx >= 0 {
// 		return cmd[idx+1:]
// 	}
// 	if len(args) == 2 {
// 		return args[1]
// 	}
// 	return ""
// }
