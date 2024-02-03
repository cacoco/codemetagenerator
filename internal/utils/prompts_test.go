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
	selection, err := MkPrompt(&stdin, &stdout, text)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "answer"
	g.Ω(*selection).Should(gomega.Equal(expected))
}
