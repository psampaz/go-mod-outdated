package mod_test

import (
	"reflect"
	"testing"

	"github.com/psampaz/go-mod-outdated/mod"
)

var mods = []mod.Module{
	mod.Module{
		Path:     "github.com/pk1/pk1",
		Main:     true,
		Indirect: false,
	},
	mod.Module{
		Path:     "github.com/pk2/pk2",
		Main:     false,
		Version:  "v1.0.0",
		Indirect: false,
	},
	mod.Module{
		Path:     "github.com/pk3/pk3",
		Main:     false,
		Version:  "v1.0.0",
		Indirect: true,
	},
	mod.Module{
		Path:     "github.com/pk4/pk4",
		Main:     false,
		Version:  "v1.0.0",
		Indirect: false,
		Update: &mod.Module{
			Version: "v1.1.0",
		},
	},
	mod.Module{
		Path:     "github.com/pk4/pk4",
		Main:     false,
		Version:  "v1.0.0",
		Indirect: true,
		Update: &mod.Module{
			Version: "v1.1.0",
		},
	},
}

func TestFilterModules(t *testing.T) {
	want := []mod.Module{
		mod.Module{
			Path:     "github.com/pk2/pk2",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
		},
		mod.Module{
			Path:     "github.com/pk3/pk3",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: true,
		},
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: true,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
	}

	if got := mod.FilterModules(mods, false, false); !reflect.DeepEqual(got, want) {
		t.Errorf("FilterModules() = %v, want %v", got, want)
	}
}

func TestFilterModulesHasUpdate(t *testing.T) {
	want := []mod.Module{
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: true,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
	}

	if got := mod.FilterModules(mods, true, false); !reflect.DeepEqual(got, want) {
		t.Errorf("FilterModules() = %v, want %v", got, want)
	}
}

func TestFilterModulesIsDirect(t *testing.T) {
	want := []mod.Module{
		mod.Module{
			Path:     "github.com/pk2/pk2",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
		},
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
	}

	if got := mod.FilterModules(mods, false, true); !reflect.DeepEqual(got, want) {
		t.Errorf("FilterModules() = %v, want %v", got, want)
	}
}

func TestFilterModulesHasUpdateIsDirect(t *testing.T) {
	want := []mod.Module{
		mod.Module{
			Path:     "github.com/pk4/pk4",
			Main:     false,
			Version:  "v1.0.0",
			Indirect: false,
			Update: &mod.Module{
				Version: "v1.1.0",
			},
		},
	}

	if got := mod.FilterModules(mods, true, true); !reflect.DeepEqual(got, want) {
		t.Errorf("FilterModules() = %v, want %v", got, want)
	}
}
