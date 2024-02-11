package utils

import (
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
	"time"
)

type Predicate[T any] func(T) bool

func Filter[T any](input []T, predicate Predicate[T]) []T {
	out := []T{}
	for _, elem := range input {
		if predicate(elem) {
			out = append(out, elem)
		}
	}
	return out
}

func ValidNumber(a any) error {
	switch a.(type) {
	case int:
		return nil
	case float64:
		return nil
	default:
		_, err := strconv.Atoi(a.(string))
		return err
	}
}

func ValidUrl(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url: %s", str)
	}
	return nil
}

func ValidEmailAddress(address string) error {
	_, err := mail.ParseAddress(address)
	if err != nil {
		return fmt.Errorf("invalid email address: %s", address)
	}
	return nil
}

func ValidBoolean(str string) error {
	if str != "true" && str != "false" {
		return fmt.Errorf("invalid boolean: %s", str)
	}
	return nil
}

func ValidISO8601Date(str string) error {
	_, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return fmt.Errorf("invalid ISO8601 date: %s", str)
	}
	return nil
}
