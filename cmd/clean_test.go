package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/cacoco/codemetagenerator/internal/utils"
	"github.com/spf13/cobra"
)

func Test_ExecuteCleanCmd1(t *testing.T) {
	temp := t.TempDir()
	// setup
	os.Mkdir(utils.GetHomeDir(temp), 0755)

	writer := &utils.TestWriter{}

	new := &cobra.Command{Use: "new", RunE: func(cmd *cobra.Command, args []string) error {
		return clean(temp, writer)
	},
	}
	buf := bytes.NewBufferString("")
	new.SetOut(buf)
	new.SetErr(buf)
	new.SetArgs([]string{})

	err := new.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check directory no longer exists
	if _, err := os.Stat(utils.GetHomeDir(temp)); err == nil {
		t.Errorf("Directory should have been removed")
	}
}

func Test_ExecuteCleanCmd2(t *testing.T) {
	// shouldn't fail when directory doesn't exist
	temp := t.TempDir()

	writer := &utils.TestWriter{}

	new := &cobra.Command{Use: "new", RunE: func(cmd *cobra.Command, args []string) error {
		return clean(temp, writer)
	},
	}
	buf := bytes.NewBufferString("")
	new.SetOut(buf)
	new.SetErr(buf)
	new.SetArgs([]string{})

	err := new.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// check directory STILL no longer exists
	if _, err := os.Stat(utils.GetHomeDir(temp)); err == nil {
		t.Errorf("Directory should not exist")
	}
}
