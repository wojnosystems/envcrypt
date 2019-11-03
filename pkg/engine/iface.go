package engine

import (
	"envcrypt/pkg/envfile"
)

type Enginer interface {
	Encrypt(envFile envfile.EnvUnsealed) (envSealed envfile.EnvSealed, err error)
	Decrypt(sealed envfile.EnvSealed) (envFile envfile.EnvUnsealed, err error)
}
