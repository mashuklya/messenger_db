package utils

import (
	"errors"
	"regexp"
	"strings"
)

func IsUUID(id string) bool {
	patternUUID := "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"
	ok, _ := regexp.Match(patternUUID, []byte(id))
	return ok
}

func ConcatErrors(errs []error) error {
	var sErr string
	for _, err := range errs {
		sErr = strings.Join([]string{sErr, err.Error()}, "\n")
	}
	return errors.New(sErr)
}
