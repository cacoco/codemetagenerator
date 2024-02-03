package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func Test_ExecuteAddCmd(t *testing.T) {
	add := &cobra.Command{Use: "add", RunE: addCmdRunE}
	b := bytes.NewBufferString("")
	add.SetOut(b)
	add.SetErr(b)

	err := add.Execute() // requires a sub-command argument
	if err == nil {
		t.Errorf("expected error for add command")
	}
}

func Test_ExecuteAddCmd1(t *testing.T) {
	add := &cobra.Command{Use: "add", Args: checkArgs, RunE: addCmdRunE}
	b := bytes.NewBufferString("")
	add.SetOut(b)
	add.SetErr(b)

	add.SetArgs([]string{"author", "badarg"}) // should only have the single sub-command arg

	err := add.Execute()
	if err == nil {
		t.Errorf("expected error for add command")
	}
}

func Test_ExecuteAddCmd2(t *testing.T) {
	add := &cobra.Command{Use: "add", Args: checkArgs, RunE: addCmdRunE}
	b := bytes.NewBufferString("")
	add.SetOut(b)
	add.SetErr(b)

	add.SetArgs([]string{"contributor", "badarg"}) // should only have the single sub-command arg

	err := add.Execute()
	if err == nil {
		t.Errorf("expected error for add command")
	}
}

func Test_ExecuteAddCmd3(t *testing.T) {
	add := &cobra.Command{Use: "add", Args: checkArgs, RunE: addCmdRunE}
	b := bytes.NewBufferString("")
	add.SetOut(b)
	add.SetErr(b)

	add.SetArgs([]string{"keyword"}) // should have at least one more arg

	err := add.Execute()
	if err == nil {
		t.Errorf("expected error for add command")
	}
}

func Test_ExecuteAddCmd4(t *testing.T) {
	add := &cobra.Command{Use: "add", Args: checkArgs, RunE: addCmdRunE}
	b := bytes.NewBufferString("")
	add.SetOut(b)
	add.SetErr(b)

	add.SetArgs([]string{"unrecognized"}) // not a valid sub-command

	err := add.Execute()
	if err == nil {
		t.Errorf("expected error for add command")
	}
}
