package envfile

import (
	"encoding/binary"
	"fmt"
	"github.com/wojnosystems/envcrypt/pkg/stream_io"
	"io"
)

type EnvSealed struct {
	version        uint16
	engineType     string
	engineMeta     []byte
	sealedContents []byte
}

const (
	MaximumEngineTypeLength     = 32
	MaximumEngineMetaLength     = 1024
	MaximumSealedContentsLength = 1_000_000_000
)

func NewEnvSealed(engineType string) (sealed EnvSealed, err error) {
	sealed = EnvSealed{
		version:        1,
		engineType:     engineType,
		engineMeta:     nil,
		sealedContents: nil,
	}
	err = sealed.SetEngineType(engineType)
	return
}

func (e *EnvSealed) SetEngineType(engineType string) (err error) {
	if len(engineType) > MaximumEngineTypeLength {
		return fmt.Errorf("engineType string is too long, maximum length is %d", MaximumEngineTypeLength)
	}
	e.engineType = engineType
	return nil
}

func (e *EnvSealed) SetEngineMeta(engineMeta []byte) (err error) {
	if len(engineMeta) > MaximumEngineMetaLength {
		return fmt.Errorf("engineMeta is too long, maximum length is %d", MaximumEngineMetaLength)
	}
	e.engineMeta = engineMeta
	return nil
}

func (e *EnvSealed) SetSealedEnvContents(sealedContents []byte) {
	e.sealedContents = sealedContents
}

func (e *EnvSealed) EngineType() (sealedContents string) {
	return e.engineType
}

func (e *EnvSealed) EngineMeta() (engineMeta []byte) {
	return e.engineMeta
}

func (e *EnvSealed) SealedEnvContents() (sealedContents []byte) {
	return e.sealedContents
}

func (e EnvSealed) Write(writer io.Writer) (err error) {
	err = binary.Write(writer, binary.BigEndian, e.version)
	if err != nil {
		return
	}
	err = stream_io.WriteString(writer, e.engineType)
	if err != nil {
		return
	}

	err = stream_io.WriteByteSlice(writer, e.engineMeta)
	if err != nil {
		return
	}

	err = stream_io.WriteByteSlice(writer, e.sealedContents)
	return
}

func (e *EnvSealed) Read(reader io.Reader) (err error) {
	err = binary.Read(reader, binary.BigEndian, &e.version)
	if err != nil {
		return
	}
	e.engineType, err = stream_io.ReadString(reader, MaximumEngineTypeLength)
	if err != nil {
		return
	}

	e.engineMeta = make([]byte, MaximumEngineMetaLength)
	e.engineMeta, err = stream_io.ReadIntoByteSlice(reader, e.engineMeta)
	if err != nil {
		return
	}

	e.sealedContents, err = stream_io.ReadByteSlice(reader, MaximumSealedContentsLength)
	return
}
