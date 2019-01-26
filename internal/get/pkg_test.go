// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package get

import (
	"github.com/tiborvass/gomod/internal/str"
	"reflect"
	"strings"
	"testing"
)

var foldDupTests = []struct {
	list   []string
	f1, f2 string
}{
	{str.StringList("math/rand", "math/big"), "", ""},
	{str.StringList("math", "strings"), "", ""},
	{str.StringList("strings"), "", ""},
	{str.StringList("strings", "strings"), "strings", "strings"},
	{str.StringList("Rand", "rand", "math", "math/rand", "math/Rand"), "Rand", "rand"},
}

func TestFoldDup(t *testing.T) {
	for _, tt := range foldDupTests {
		f1, f2 := str.FoldDup(tt.list)
		if f1 != tt.f1 || f2 != tt.f2 {
			t.Errorf("foldDup(%q) = %q, %q, want %q, %q", tt.list, f1, f2, tt.f1, tt.f2)
		}
	}
}

var parseMetaGoImportsTests = []struct {
	in  string
	mod ModuleMode
	out []metaImport
}{
	{
		`<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">`,
		IgnoreMod,
		[]metaImport{{"foo/bar", "git", "https://github.com/rsc/foo/bar"}},
	},
	{
		`<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">
		<meta name="go-import" content="baz/quux git http://github.com/rsc/baz/quux">`,
		IgnoreMod,
		[]metaImport{
			{"foo/bar", "git", "https://github.com/rsc/foo/bar"},
			{"baz/quux", "git", "http://github.com/rsc/baz/quux"},
		},
	},
	{
		`<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">
		<meta name="go-import" content="foo/bar mod http://github.com/rsc/baz/quux">`,
		IgnoreMod,
		[]metaImport{
			{"foo/bar", "git", "https://github.com/rsc/foo/bar"},
		},
	},
	{
		`<meta name="go-import" content="foo/bar mod http://github.com/rsc/baz/quux">
		<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">`,
		IgnoreMod,
		[]metaImport{
			{"foo/bar", "git", "https://github.com/rsc/foo/bar"},
		},
	},
	{
		`<meta name="go-import" content="foo/bar mod http://github.com/rsc/baz/quux">
		<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">`,
		PreferMod,
		[]metaImport{
			{"foo/bar", "mod", "http://github.com/rsc/baz/quux"},
		},
	},
	{
		`<head>
		<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">
		</head>`,
		IgnoreMod,
		[]metaImport{{"foo/bar", "git", "https://github.com/rsc/foo/bar"}},
	},
	{
		`<meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">
		<body>`,
		IgnoreMod,
		[]metaImport{{"foo/bar", "git", "https://github.com/rsc/foo/bar"}},
	},
	{
		`<!doctype html><meta name="go-import" content="foo/bar git https://github.com/rsc/foo/bar">`,
		IgnoreMod,
		[]metaImport{{"foo/bar", "git", "https://github.com/rsc/foo/bar"}},
	},
	{
		// XML doesn't like <div style=position:relative>.
		`<!doctype html><title>Page Not Found</title><meta name=go-import content="chitin.io/chitin git https://github.com/chitin-io/chitin"><div style=position:relative>DRAFT</div>`,
		IgnoreMod,
		[]metaImport{{"chitin.io/chitin", "git", "https://github.com/chitin-io/chitin"}},
	},
	{
		`<meta name="go-import" content="myitcv.io git https://github.com/myitcv/x">
	        <meta name="go-import" content="myitcv.io/blah2 mod https://raw.githubusercontent.com/myitcv/pubx/master">
	        `,
		IgnoreMod,
		[]metaImport{{"myitcv.io", "git", "https://github.com/myitcv/x"}},
	},
	{
		`<meta name="go-import" content="myitcv.io git https://github.com/myitcv/x">
	        <meta name="go-import" content="myitcv.io/blah2 mod https://raw.githubusercontent.com/myitcv/pubx/master">
	        `,
		PreferMod,
		[]metaImport{
			{"myitcv.io/blah2", "mod", "https://raw.githubusercontent.com/myitcv/pubx/master"},
			{"myitcv.io", "git", "https://github.com/myitcv/x"},
		},
	},
}

func TestParseMetaGoImports(t *testing.T) {
	for i, tt := range parseMetaGoImportsTests {
		out, err := parseMetaGoImports(strings.NewReader(tt.in), tt.mod)
		if err != nil {
			t.Errorf("test#%d: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf("test#%d:\n\thave %q\n\twant %q", i, out, tt.out)
		}
	}
}
