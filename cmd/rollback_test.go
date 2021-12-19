package cmd

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/roots/trellis-cli/trellis"
)

func TestRollbackRunValidations(t *testing.T) {
	defer trellis.LoadFixtureProject(t)()

	cases := []struct {
		name            string
		projectDetected bool
		args            []string
		out             string
		code            int
	}{
		{
			"no_project",
			false,
			nil,
			"No Trellis project detected",
			1,
		},
		{
			"no_args",
			true,
			nil,
			"Usage: trellis",
			1,
		},
		{
			"invalid_env",
			true,
			[]string{"foo"},
			"Error: foo is not a valid environment",
			1,
		},
		{
			"invalid_env_with_site",
			true,
			[]string{"foo", "example.com"},
			"Error: foo is not a valid environment",
			1,
		},
		{
			"invalid_site",
			true,
			[]string{"development", "nosite"},
			"Error: nosite is not a valid site",
			1,
		},
		{
			"too_many_args",
			true,
			[]string{"development", "site", "foo"},
			"Error: too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		ui := cli.NewMockUi()
		mockProject := &MockProject{tc.projectDetected}
		trellis := trellis.NewTrellis(mockProject)
		rollbackCommand := NewRollbackCommand(ui, trellis)

		code := rollbackCommand.Run(tc.args)

		if code != tc.code {
			t.Errorf("expected code %d to be %d", code, tc.code)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		if !strings.Contains(combined, tc.out) {
			t.Errorf("expected output %q to contain %q", combined, tc.out)
		}
	}
}

func TestRollbackRun(t *testing.T) {
	defer trellis.LoadFixtureProject(t)()
	project := &trellis.Project{}
	trellis := trellis.NewTrellis(project)

	defer MockExec(t)()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"default",
			[]string{"development", "example.com"},
			"ansible-playbook rollback.yml -e env=development site=example.com",
			0,
		},
		{
			"site_not_needed_in_defaut_case",
			[]string{"development"},
			"ansible-playbook rollback.yml -e env=development site=example.com",
			0,
		},
		{
			"with_release_flag",
			[]string{"--release=123", "development", "example.com"},
			"ansible-playbook rollback.yml -e env=development site=example.com release=123",
			0,
		},
	}

	for _, tc := range cases {
		ui := cli.NewMockUi()
		rollbackCommand := NewRollbackCommand(ui, trellis)

		code := rollbackCommand.Run(tc.args)

		if code != tc.code {
			t.Errorf("expected code %d to be %d", code, tc.code)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		if !strings.Contains(combined, tc.out) {
			t.Errorf("expected output %q to contain %q", combined, tc.out)
		}
	}
}
