package utils

import (
	"fmt"
	"net/mail"
	"net/url"
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
