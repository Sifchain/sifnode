package utils

import "github.com/pkg/errors"

func ParseStringToBool(s string) (bool, error) {
	if s == "true" || s == "True" || s == "TRUE" {
		return true, nil
	}
	if s == "false" || s == "False" || s == "FALSE" {
		return false, nil
	}
	return false, errors.New("Can only accept true or false")
}
