package cmd

import (
	"testing"

	"github.com/cacoco/codemetagenerator/internal/model"
	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/onsi/gomega"
)

func TestNewAuthor(t *testing.T) {
	g := gomega.NewWithT(t)

	var stack utils.Stack[string]
	stack.Push("id\n")
	stack.Push("person@email.org\n")
	stack.Push("familyName\n")
	stack.Push("givenName\n")
	stack.Push("\n") // enter to select the first option
	reader := utils.TestReader{In: utils.TestStdin{Data: stack}}
	writer := utils.TestWriter{}

	author, err := newAuthor(&reader, &writer)
	if author == nil {
		t.Errorf("Expected author to not be nil")
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := map[string]any{
		model.Type:       model.PersonType,
		model.FamilyName: "familyName",
		model.GivenName:  "givenName",
		model.Email:      "person@email.org",
		model.Id:         "id",
	}
	g.Ω(*author).Should(gomega.Equal(expected))
}
