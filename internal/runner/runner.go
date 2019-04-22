package runner

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/psampaz/go-mod-outdated/internal/mod"

	"github.com/olekukonko/tablewriter"
)

// Run converts the the json output of go list -u -m -json all to table format
func Run(in io.Reader, out io.Writer, update, direct bool) error {
	var modules []mod.Module
	dec := json.NewDecoder(in)

	for {
		var v mod.Module
		err := dec.Decode(&v)

		if err != nil {
			if err == io.EOF {
				renderTable(out, mod.FilterModules(modules, update, direct))
				return nil
			}
			return err
		}
		modules = append(modules, v)
	}
}

func renderTable(writer io.Writer, modules []mod.Module) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Module", "Version", "New Version", "Direct", "Valid Timestamps"})
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
}
