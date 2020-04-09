// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// notgo.mod init

package modcmd

import (
	"github.com/tiborvass/gomod/internal/cmd/go/base"
	"github.com/tiborvass/gomod/internal/cmd/go/modload"
	"github.com/tiborvass/gomod/internal/cmd/go/work"
	"os"
	"strings"
)

var cmdInit = &base.Command{
	UsageLine: "notgo.mod init [module]",
	Short:     "initialize new module in current directory",
	Long: `
Init initializes and writes a new notgo.mod to the current directory,
in effect creating a new module rooted at the current directory.
The file notgo.mod must not already exist.
If possible, init will guess the module path from import comments
(see 'go help importpath') or from version control configuration.
To override this guess, supply the module path as an argument.
	`,
	Run: runInit,
}

func init() {
	work.AddModCommonFlags(cmdInit)
}

func runInit(cmd *base.Command, args []string) {
	modload.CmdModInit = true
	if len(args) > 1 {
		base.Fatalf("notgo.mod init: too many arguments")
	}
	if len(args) == 1 {
		modload.CmdModModule = args[0]
	}
	if os.Getenv("GO111MODULE") == "off" {
		base.Fatalf("notgo.mod init: modules disabled by GO111MODULE=off; see 'go help modules'")
	}
	modFilePath := modload.ModFilePath()
	if _, err := os.Stat(modFilePath); err == nil {
		base.Fatalf("notgo.mod init: notgo.mod already exists")
	}
	if strings.Contains(modload.CmdModModule, "@") {
		base.Fatalf("notgo.mod init: module path must not contain '@'")
	}
	modload.InitMod() // does all the hard work
}