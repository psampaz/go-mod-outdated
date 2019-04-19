package mod

// Module holds information about a specifc module listed by go list
type Module struct {
	Path     string `json:"Path"`
	Main     bool   `json:"Main"`
	Version  string `json:"Version"`
	Indirect bool   `json:"Indirect"`
	Update   Update `json:"Update,omitempty"`
}

// Update holds information about the updates of a specific module
type Update struct {
	Version string `json:"Version"`
}

// FilterModules filters the list of modules provided by the go list command
func FilterModules(modules []Module, hasUpdate, isDirect bool) []Module {
	out := make([]Module, 0)
	for k := range modules {

		if modules[k].Main {
			continue
		}

		if hasUpdate && modules[k].Update.Version == "" {
			continue
		}

		if isDirect && modules[k].Indirect {
			continue
		}

		out = append(out, modules[k])
	}

	return out
}
