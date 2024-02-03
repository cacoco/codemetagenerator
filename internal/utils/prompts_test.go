package utils

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestMkPrompt(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack Stack[string]
	stack.Push("answer\n")

	reader := TestReader{In: TestStdin{Data: stack}}
	stdin := reader.Stdin()

	writer := TestWriter{}
	stdout := writer.Stdout()

	text := "test"
	selection, err := MkPrompt(&stdin, &stdout, text, Nop)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "answer"
	g.Ω(*selection).Should(gomega.Equal(expected))
}

func TestMkPromptInvalidURL(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack Stack[string]
	stack.Push("NotaURL\n")

	reader := TestReader{In: TestStdin{Data: stack}}
	stdin := reader.Stdin()

	writer := TestWriter{}
	stdout := writer.Stdout()

	text := "test"
	selection, err := MkPrompt(&stdin, &stdout, text, ValidUrl)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	g.Expect(err).ToNot(gomega.BeNil())
	g.Ω(err.Error()).Should(gomega.Equal("invalid url: NotaURL"))
	// errored, selection should be nil
	g.Expect(selection).To(gomega.BeNil())
}

func TestMkPromptValidURL(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack Stack[string]
	stack.Push("https://google.com\n")

	reader := TestReader{In: TestStdin{Data: stack}}
	stdin := reader.Stdin()

	writer := TestWriter{}
	stdout := writer.Stdout()

	text := "test"
	selection, err := MkPrompt(&stdin, &stdout, text, ValidUrl)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(selection).ToNot(gomega.BeNil())
	g.Ω(*selection).Should(gomega.Equal("https://google.com"))
}

func TestMkPromptInvalidEmailAddress(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack Stack[string]
	stack.Push("NotanEmail\n")

	reader := TestReader{In: TestStdin{Data: stack}}
	stdin := reader.Stdin()

	writer := TestWriter{}
	stdout := writer.Stdout()

	text := "test"
	selection, err := MkPrompt(&stdin, &stdout, text, ValidEmailAddress)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	g.Expect(err).ToNot(gomega.BeNil())
	g.Ω(err.Error()).Should(gomega.Equal("invalid email address: NotanEmail"))
	// errored, selection should be nil
	g.Expect(selection).To(gomega.BeNil())
}

func TestMkPromptValidEmailAddress(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack Stack[string]
	stack.Push("person@email.org\n")

	reader := TestReader{In: TestStdin{Data: stack}}
	stdin := reader.Stdin()

	writer := TestWriter{}
	stdout := writer.Stdout()

	text := "test"
	selection, err := MkPrompt(&stdin, &stdout, text, ValidEmailAddress)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(selection).ToNot(gomega.BeNil())
	g.Ω(*selection).Should(gomega.Equal("person@email.org"))
}
