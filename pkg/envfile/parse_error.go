package envfile

import "fmt"

type parseError struct {
	lineNumber int
	reason     string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("unable to parse environment file, %s on line: %d", e.reason, e.lineNumber)
}

func newParseError(lineNumber int, reason string) *parseError {
	return &parseError{
		lineNumber: lineNumber,
		reason:     reason,
	}
}
