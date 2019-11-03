package envfile

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

type EnvUnsealed map[string]string

func FromUnsealedFile(reader io.Reader) (envs EnvUnsealed, err error) {
	envs = make(EnvUnsealed)
	scanner := bufio.NewScanner(reader)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			return envs, newParseError(lineNumber, "no equal (=) sign")
		}
		key := parts[0]
		value := parts[1]
		envs[key] = value
		lineNumber++
	}
	return
}

func (e EnvUnsealed) ToUnsealedFile(writer io.Writer) (bytesWritten int, err error) {
	var bytes int
	orderedKeys := make([]string, 0, len(e))
	for key := range e {
		orderedKeys = append(orderedKeys, key)
	}
	sort.Strings(orderedKeys)
	for _, key := range orderedKeys {
		bytes, err = writer.Write([]byte(key))
		if err != nil {
			return
		}
		bytesWritten += bytes

		bytes, err = writer.Write([]byte("="))
		if err != nil {
			return
		}
		bytesWritten += bytes

		bytes, err = writer.Write([]byte(e[key]))
		if err != nil {
			return
		}
		bytesWritten += bytes

		bytes, err = writer.Write([]byte("\n"))
		if err != nil {
			return
		}
		bytesWritten += bytes
	}
	return
}
