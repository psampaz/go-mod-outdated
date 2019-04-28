package runner

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/psampaz/go-mod-outdated/internal/mod"

	"github.com/olekukonko/tablewriter"
)

// Run converts the the json output of go list -u -m -json all to table format
func Run(in io.Reader, out io.Writer, update, direct, exitWithNonZero bool) error {
	var modules []mod.Module
	dec := json.NewDecoder(in)

	for {
		var v mod.Module
		err := dec.Decode(&v)

		if err != nil {
			if err == io.EOF {
				found := renderTable(out, mod.FilterModules(modules, update, direct))
				if found && exitWithNonZero {
					os.Exit(1)
				}
				return nil
			}
			return err
		}
		modules = append(modules, v)
	}
}

func renderTable(writer io.Writer, modules []mod.Module) bool {
	var found bool
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Module", "Version", "New Version", "Direct", "Valid Timestamps"})
	for k := range modules {
		if modules[k].NewVersion() != "" {
			found = true
		}
		table.Append([]string{
			modules[k].Path,
			modules[k].CurrentVersion(),
			modules[k].NewVersion(),
			strconv.FormatBool(!modules[k].Indirect),
			strconv.FormatBool(!modules[k].InvalidTimestamp()),
		})
	}
	table.Render()
	return found
}
