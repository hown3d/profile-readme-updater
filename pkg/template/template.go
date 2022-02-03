package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/hown3d/profile-readme-updater/pkg/github"
)

func Render(out io.Writer, filepath string, event *github.Infos) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file on %v: %w", filepath, err)
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t := template.Must(template.New("readme").Funcs(funcMap).Parse(string(data)))

	err = t.Execute(out, event)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}
