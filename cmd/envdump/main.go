package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for _, keyAndValue := range os.Environ() {
		parts := strings.SplitN(keyAndValue, "=", 2)
		key := parts[0]
		value := parts[1]
		fmt.Printf("Env: %s=%s\n", key, value)
	}
}
