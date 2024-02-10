package utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/manifoldco/promptui"
)

func Nop(s string) error {
	return nil
}

func MkPrompt(stdin *io.ReadCloser, stdout *io.WriteCloser, text string, validate func(string) error) (*string, error) {
	prompt := promptui.Prompt{
		Label:  text,
		Stdin:  *stdin,
		Stdout: *stdout,
	}

	result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	// the promptui library validation fails in odd ways, so we manually validate here
	err = validate(result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewPersonOrOrganizationPrompt(reader *Reader, writer *Writer, label string) (*map[string]any, error) {
	stdin := (*reader).Stdin()
	stdout := (*writer).Stdout()

	options := []model.MenuOption{
		{Name: "Person", Type: "person"},
		{Name: "Organization", Type: "organization"},
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "➞ {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: fmt.Sprintf(`{{ "Selected %s type:" | faint}} {{ .Name | faint }}`, strings.ToLower(label)),
		Details: fmt.Sprintf(`--------- %s ----------
{{ "Name:" | faint }}	{{ .Name }}`, label),
	}

	prompt := promptui.Select{
		Label:     "Please enter a " + label + " type:",
		Items:     options,
		Templates: templates,
		Size:      2,
		Searcher:  nil,
		Stdin:     stdin,
		Stdout:    stdout,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	keyType := options[i].Type
	switch keyType {
	case "person":
		givenName, err := MkPrompt(&stdin, &stdout, "Enter the given (first) name of the person", Nop)
		if err != nil {
			return nil, err
		}
		familyName, err := MkPrompt(&stdin, &stdout, "Enter the family (last) name of the person", Nop)
		if err != nil {
			return nil, err
		}
		email, err := MkPrompt(&stdin, &stdout, "Enter the email address of the person", ValidEmailAddress)
		if err != nil {
			return nil, err
		}
		id, err := MkPrompt(&stdin, &stdout, "Enter the identifier of the person (see: https://orcid.org)", Nop)
		if err != nil {
			return nil, err
		}
		return model.NewPerson(givenName, familyName, email, id), nil
	case "organization":
		name, err := MkPrompt(&stdin, &stdout, "Enter the name of the organization", Nop)
		if err != nil {
			return nil, err
		}
		url, err := MkPrompt(&stdin, &stdout, "Enter the URL of the organization", ValidUrl)
		if err != nil {
			return nil, err
		}
		id, err := MkPrompt(&stdin, &stdout, "Enter the identifier of the organization (see: https://orcid.org)", Nop)
		if err != nil {
			return nil, err
		}
		return model.NewOrganization(name, url, id), nil
	default:
		return nil, fmt.Errorf("Invalid selection: " + keyType)
	}
}
