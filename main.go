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
				Usage:     "Add, replace or remove a configured secret. To remove a secret, set it's secret token to empty.",
				UsageText: "totpgen set <secret-name> <secret-token>\ntotpgen set <secret-name> ''",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() != 2 {
						cli.ShowSubcommandHelpAndExit(ctx, 1)
					}
					ring, err := openKeyring()
					if err != nil {
						return err
					}
					secret := ctx.Args().Get(1)
					if secret == "" {
						if err := ring.Remove(ctx.Args().Get(0)); err != nil {
							return fmt.Errorf("could not remove secret from keyring: %w", err)
						}
					}

					err = ring.Set(keyring.Item{
						Key:  ctx.Args().Get(0),
						Data: []byte(secret),
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
			{
				Name:      "backends",
				Usage:     "show available keyring backends",
				UsageText: "totpgen backends",
				Action: func(ctx *cli.Context) error {
					for _, b := range keyring.AvailableBackends() {
						fmt.Println(b)
					}
					return nil
				},
			},
			// better do not offer export for security reasons.
			// A later version might allow this, but maybe encrypt the secrets with a password or PK.
			// {
			// 	Name:      "export",
			// 	Usage:     "export configured secrets",
			// 	UsageText: "totpgen export",
			// 	Action: func(ctx *cli.Context) error {
			// 		ring, err := openKeyring()
			// 		if err != nil {
			// 			return err
			// 		}
			// 		kk, err := ring.Keys()
			// 		if err != nil {
			// 			return fmt.Errorf("could not list keys from keyring: %w", err)
			// 		}
			// 		sort.Strings(kk)
			// 		for _, k := range kk {
			// 			e, err := ring.Get(k)
			// 			if err != nil {
			// 				return fmt.Errorf("could not get secret for '%s': %w", k, err)
			// 			}
			// 			fmt.Printf(`totpgen set "%s" "%s"`, k, e.Data)
			// 			fmt.Println()
			// 		}
			// 		return nil
			// 	},
			// },
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
		ServiceName:  ringName,
		KeychainName: os.Getenv("TOTPGEN_KEYCHAIN_NAME"),
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
