package mod

import "time"

// Module holds information about a specifc module listed by go list

type Module struct {
	Path      string        `json:",omitempty"` // module path
	Version   string        `json:",omitempty"` // module version
	Versions  []string      `json:",omitempty"` // available module versions
	Replace   *Module `json:",omitempty"` // replaced by this module
	Time      *time.Time    `json:",omitempty"` // time version was created
	Update    *Module `json:",omitempty"` // available update (with -u)
	Main      bool          `json:",omitempty"` // is this the main module?
	Indirect  bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string        `json:",omitempty"` // path to go.mod file describing module, if any
	Error     *ModuleError  `json:",omitempty"` // error loading module
	GoVersion string        `json:",omitempty"` // go version used in module
}

type ModuleError struct {
	Err string // error text
}
// FilterModules filters the list of modules provided by the go list command
func FilterModules(modules []Module, hasUpdate, isDirect bool) []Module {
	out := make([]Module, 0)
	for k := range modules {

		if modules[k].Main {
			continue
		}

		if hasUpdate && modules[k].Update == nil {
			continue
		}

		if isDirect && modules[k].Indirect {
			continue
		}

		out = append(out, modules[k])
	}

	return out
}
