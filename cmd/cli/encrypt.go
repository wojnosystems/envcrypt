package main

import (
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli"
	"github.com/wojnosystems/envcrypt/pkg/engine"
	"github.com/wojnosystems/envcrypt/pkg/envfile"
	"io"
)

func encryptCommand(engineType string) func(c *cli.Context) (err error) {
	switch engineType {
	case "aes256":
		return encryptAes256Command
	default:
		return func(c *cli.Context) (err error) {
			return fmt.Errorf("invalid engine type")
		}
	}
}

func encryptAes256Command(c *cli.Context) (err error) {
	key, variables, outputStream, err := inputsForEncrypt(c)
	if err != nil {
		return
	}
	defer func() {
		_ = outputStream.Close()
	}()

	sealed, err := encryptAES256(key, variables)
	if err != nil {
		return
	}

	err = sealed.Write(outputStream)
	_ = outputStream.Close()
	return
}

func inputsForEncrypt(c *cli.Context) (key []byte, env envfile.EnvUnsealed, outputStream io.WriteCloser, err error) {
	var inputStream io.ReadCloser
	key, err = base64.StdEncoding.DecodeString(c.String("keyBase64"))
	if err != nil {
		return nil, env, nil, err
	}
	inputStream, outputStream, err = prepareStreamsForEncrypt(c.String("in"), c.String("out"))
	if err != nil {
		return
	}
	defer func() {
		_ = inputStream.Close()
	}()
	env, err = envfile.FromUnsealedFile(inputStream)
	return
}

func encryptAES256(key []byte, variables envfile.EnvUnsealed) (sealed envfile.EnvSealed, err error) {
	eng, err := engine.NewAes256(key)
	if err != nil {
		return
	}
	sealed, err = eng.Encrypt(variables)
	return
}
