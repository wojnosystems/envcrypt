package main

import (
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli"
	"github.com/wojnosystems/envcrypt/pkg/engine"
	"github.com/wojnosystems/envcrypt/pkg/envfile"
	"os"
	"syscall"
)

func decryptCommand(engineType string) func(c *cli.Context) (err error) {
	switch engineType {
	case "aes256":
		return decryptToEnv
	default:
		return func(c *cli.Context) (err error) {
			return fmt.Errorf("invalid engine type")
		}
	}
}

func decryptToEnv(c *cli.Context) (err error) {
	key, envs, err := inputsForDecrypt(c)
	if err != nil {
		return
	}

	unsealed, err := decryptAES256(key, envs)
	if err != nil {
		return
	}

	// override the key environment variable
	err = os.Unsetenv("KEY")
	if err != nil {
		return
	}

	// Set the environment variables from the now unsealed files
	for key, value := range unsealed {
		err = os.Setenv(key, value)
		if err != nil {
			return
		}
	}

	// execute the thing!
	err = syscall.Exec(c.String("exec"), c.StringSlice("execArg"), os.Environ())
	return
}

func inputsForDecrypt(c *cli.Context) (key []byte, sealedEnvs []envfile.EnvSealed, err error) {
	key, err = base64.StdEncoding.DecodeString(c.String("keyBase64"))
	if err != nil {
		return
	}

	sealedEnvs = make([]envfile.EnvSealed, 0, len(c.StringSlice("in")))
	for _, filePath := range c.StringSlice("in") {
		var file *os.File
		file, err = os.Open(filePath)
		if err != nil {
			return
		}
		sealed := envfile.EnvSealed{}
		err = sealed.Read(file)
		_ = file.Close()
		if err != nil {
			return
		}
		sealedEnvs = append(sealedEnvs, sealed)
	}
	return
}

func decryptAES256(key []byte, envs []envfile.EnvSealed) (unsealed envfile.EnvUnsealed, err error) {
	eng, err := engine.NewAes256(key)
	if err != nil {
		return
	}
	unsealed = make(envfile.EnvUnsealed)
	for _, sealedEnvs := range envs {
		var unsealedComponent envfile.EnvUnsealed
		unsealedComponent, err = eng.Decrypt(sealedEnvs)
		if err != nil {
			return
		}
		for key, value := range unsealedComponent {
			unsealed[key] = value
		}
	}
	return
}
