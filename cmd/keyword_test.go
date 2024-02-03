package cmd

import (
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
)

func TestNewKeyword(t *testing.T) {
	g := gomega.NewWithT(t)

	current := map[string]any{
		model.Keywords: []string{"one", "two"},
	}
	args := []string{"three", "four", "five"}

	writer := &utils.TestWriter{}

	keywords, err := newKeyword(writer, current, args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"one", "two", "three", "four", "five"}
	g.Ω(keywords).Should(gomega.Equal(expected))
}
