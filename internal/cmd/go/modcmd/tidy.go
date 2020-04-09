// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// notgo.mod tidy

package modcmd

import (
	"github.com/tiborvass/gomod/internal/cmd/go/base"
	"github.com/tiborvass/gomod/internal/cmd/go/cfg"
	"github.com/tiborvass/gomod/internal/cmd/go/modfetch"
	"github.com/tiborvass/gomod/internal/cmd/go/modload"
	"github.com/tiborvass/gomod/internal/cmd/go/work"

	"golang.org/x/mod/module"
)

var cmdTidy = &base.Command{
	UsageLine: "notgo.mod tidy [-v]",
	Short:     "add missing and remove unused modules",
	Long: `
Tidy makes sure notgo.mod matches the source code in the module.
It adds any missing modules necessary to build the current module's
packages and dependencies, and it removes unused modules that
don't provide any relevant packages. It also adds any missing entries
to go.sum and removes any unnecessary ones.

The -v flag causes tidy to print information about removed modules
to standard error.
	`,
}

func init() {
	cmdTidy.Run = runTidy // break init cycle
	cmdTidy.Flag.BoolVar(&cfg.BuildV, "v", false, "")
	work.AddModCommonFlags(cmdTidy)
}

func runTidy(cmd *base.Command, args []string) {
	if len(args) > 0 {
		base.Fatalf("notgo.mod tidy: no arguments allowed")
	}

	modload.LoadALL()
	modload.TidyBuildList()
	modTidyGoSum() // updates memory copy; WriteGoMod on next line flushes it out
	modload.WriteGoMod()
}

// modTidyGoSum resets the go.sum file content
// to be exactly what's needed for the current notgo.mod.
func modTidyGoSum() {
	// Assuming go.sum already has at least enough from the successful load,
	// we only have to tell modfetch what needs keeping.
	reqs := modload.Reqs()
	keep := make(map[module.Version]bool)
	replaced := make(map[module.Version]bool)
	var walk func(module.Version)
	walk = func(m module.Version) {
		// If we build using a replacement module, keep the sum for the replacement,
		// since that's the code we'll actually use during a build.
		//
		// TODO(golang.org/issue/29182): Perhaps we should keep both sums, and the
		// sums for both sets of transitive requirements.
		r := modload.Replacement(m)
		if r.Path == "" {
			keep[m] = true
		} else {
			keep[r] = true
			replaced[m] = true
		}
		list, _ := reqs.Required(m)
		for _, r := range list {
			if !keep[r] && !replaced[r] {
				walk(r)
			}
		}
	}
	walk(modload.Target)
	modfetch.TrimGoSum(keep)
}
