package cmd

import (
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
)

func TestNewContributor(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack utils.Stack[string]
	stack.Push("id\n")
	stack.Push("https://url.com\n")
	stack.Push("name\n")
	stack.Push("j\n") // down arrow to second option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}

	writer := utils.TestWriter{}

	contributor, err := newContributor(&reader, &writer)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := map[string]any{
		model.Type: model.OrganizationType,
		model.Name: "name",
		model.URL:  "https://url.com",
		model.Id:   "id",
	}
	g.Ω(*contributor).Should(gomega.Equal(expected))
}
