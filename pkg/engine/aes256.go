package engine

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"envcrypt/pkg/envfile"
	"fmt"
	"io"
)

type aes256 struct {
	key cipher.Block
}

const (
	EngineTypeAES256 = "aes256"
)

func NewAes256(key []byte) (engine Enginer, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}
	engine = &aes256{
		key: block,
	}
	return
}

func (a *aes256) Encrypt(file envfile.EnvUnsealed) (envSealed envfile.EnvSealed, err error) {
	envSealed, err = envfile.NewEnvSealed(EngineTypeAES256)
	if err != nil {
		return
	}

	unsealed := bytes.Buffer{}
	_, err = file.ToUnsealedFile(&unsealed)
	if err != nil {
		return
	}

	var encrypter cipher.AEAD
	encrypter, err = cipher.NewGCM(a.key)
	if err != nil {
		return
	}

	var iv []byte
	iv, err = generateIV(encrypter.NonceSize())
	if err != nil {
		return
	}

	err = envSealed.SetEngineMeta(iv)
	if err != nil {
		return
	}

	sealed := make([]byte, 0, unsealed.Len())
	sealed = encrypter.Seal(sealed, iv, unsealed.Bytes(), nil)
	envSealed.SetSealedEnvContents(sealed)
	return
}

func generateIV(nonceSizeByte int) (iv []byte, err error) {
	iv = make([]byte, nonceSizeByte)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	return
}

func (a *aes256) Decrypt(sealed envfile.EnvSealed) (envFile envfile.EnvUnsealed, err error) {
	if EngineTypeAES256 != sealed.EngineType() {
		err = fmt.Errorf("wrong engine type, expected: '%s' got: '%s'", EngineTypeAES256, sealed.EngineType())
		return
	}

	var decrypter cipher.AEAD
	decrypter, err = cipher.NewGCM(a.key)
	if err != nil {
		return
	}

	iv := sealed.EngineMeta()
	if len(iv) != decrypter.NonceSize() {
		err = fmt.Errorf("wrong iv size, expected: '%d' got: '%d'", decrypter.NonceSize(), len(iv))
		return
	}

	unsealed := make([]byte, 0, len(sealed.SealedEnvContents()))
	unsealed, err = decrypter.Open(unsealed, iv, sealed.SealedEnvContents(), nil)
	if err != nil {
		return
	}

	unSealedBuffer := bytes.NewBuffer(unsealed)
	envFile, err = envfile.FromUnsealedFile(unSealedBuffer)
	return
}
