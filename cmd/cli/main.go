package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "envcrypt"
	app.Usage = "Encrypt secrets, commit to the repo, deploy them in production safely"
	app.Version = "v1.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "encrypt",
			Usage: "encrypt a file",
			Subcommands: []cli.Command{
				{
					Name: "aes256",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "keyBase64",
							EnvVar:   "KEY_BASE64",
							Required: true,
						},
						cli.StringFlag{
							Name:     "in",
							Required: true,
						},
						cli.StringFlag{
							Name:     "out",
							Required: true,
						},
					},
					Action: encryptCommand("aes256"),
				},
			},
		},
		{
			Name:  "decrypt-to-env",
			Usage: "decrypt the files, run the provided executable with the secrets as environment variables",
			Subcommands: []cli.Command{
				{
					Name: "aes256",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "keyBase64",
							EnvVar:   "KEY_BASE64",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "in",
							Required: true,
						},
						cli.StringFlag{
							Name:     "exec",
							EnvVar:   "EXEC",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "execArg",
							Required: false,
						},
					},
					Action: decryptCommand("aes256"),
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
