package codemetagenerator

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func Test_ExecuteEditCmd(t *testing.T) {
	edit := &cobra.Command{Use: "edit", Args: cobra.ExactArgs(2), RunE: EditCmdRunE}
	b := bytes.NewBufferString("")
	edit.SetOut(b)
	edit.SetErr(b)

	err := edit.Execute() // required a key and value arguments
	if err == nil {
		t.Errorf("expected error for add command")
	}
}
