// Package runner is responsible for running the command and rendering the output
package runner

import (
	"encoding/json"
	"html/template"
	"io"
	"os"
	"strconv"

	"github.com/psampaz/go-mod-outdated/internal/mod"

	"github.com/olekukonko/tablewriter"
)

// OsExit is use here in order to simplify testing
var OsExit = os.Exit

// Run converts the the json output of go list -u -m -json all to table format
func Run(in io.Reader, out io.Writer, update, direct, exitWithNonZero bool, style string) error {
	var modules []mod.Module

	dec := json.NewDecoder(in)

	for {
		var v mod.Module
		err := dec.Decode(&v)

		if err != nil {
			if err == io.EOF {
				filteredModules := mod.FilterModules(modules, update, direct)
				tableErr := renderTable(out, filteredModules, style)
				if tableErr != nil {
					return tableErr
				}

				if hasOutdated(filteredModules) && exitWithNonZero {
					OsExit(1)
				}

				return nil
			}

			return err
		}

		modules = append(modules, v)
	}
}

func hasOutdated(filteredModules []mod.Module) bool {
	for m := range filteredModules {
		if filteredModules[m].HasUpdate() {
			return true
		}
	}

	return false
}

func renderTable(writer io.Writer, modules []mod.Module, style string) error {
	if style == "html" {
		return RenderHTMLTable(writer, modules)
	}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Module", "Version", "New Version", "Direct", "Valid Timestamps"})

	// Render table as markdown
	if style == "markdown" {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	}

	for k := range modules {
		table.Append([]string{
			modules[k].Path,
			modules[k].CurrentVersion(),
			modules[k].NewVersion(),
			strconv.FormatBool(!modules[k].Indirect),
			strconv.FormatBool(!modules[k].InvalidTimestamp()),
		})
	}

	table.Render()
	return nil
}

func RenderHTMLTable(writer io.Writer, modules []mod.Module) error {
	type htmlTemplateData struct {
		Path           string
		CurrentVersion string
		NewVersion     string
		Direct         bool
		ValidTimestamp bool
	}

	tableTemplate, err := template.New("dependencies").Parse(`<table><tr><th>Module</th><th>Version</th><th>New Version</th><th>Direct</th><th>Valid Timestamps</th>{{range .}}<tr><td>{{.Path}}</td><td>{{.CurrentVersion}}</td><td>{{.NewVersion}}</td><td>{{.Direct}}</td><td>{{.ValidTimestamp}}</td></tr>{{end}}</table>`)
	if err != nil {
		return err
	}
	data := make([]htmlTemplateData, len(modules), len(modules))
	for i, module := range modules {
		data[i] = htmlTemplateData{
			Path:           module.Path,
			CurrentVersion: module.CurrentVersion(),
			NewVersion:     module.NewVersion(),
			Direct:         !module.Indirect,
			ValidTimestamp: !module.InvalidTimestamp(),
		}
	}
	err = tableTemplate.Execute(writer, data)
	if err != nil {
		return err
	}
	return nil
}
