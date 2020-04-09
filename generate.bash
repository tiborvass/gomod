#!/bin/bash -e

rm -rf internal
mkdir -p internal/cmd/go
: ${GOROOT=$(go env GOROOT)}
cp -rf $GOROOT/src/cmd/go/internal/* internal/cmd/go/
cp -rf $GOROOT/src/cmd/internal/* internal/cmd/
cp -rf $GOROOT/src/internal/* internal/

cp $GOROOT/src/cmd/go/main.go .
sed -E -i'' '/func init\(\) \{/,/^\}/d' main.go

sed -E -i'' '/go11tag/d' main.go
sed -E -i''  '/go:generate/d' main.go

cat >> main.go <<EOF

func init() {
	base.Usage = mainUsage
	base.Go.Commands = []*base.Command{
		modcmd.CmdMod,
	}
}
EOF
gofmt -w -s main.go
goimports -w main.go

find . -name '*_test.go' -exec rm {} \;

for x in 'cmd/go/internal' 'cmd/internal' 'internal'; do for f in $(find . -name '*.go' -exec grep -H "\"$x/" {} \; | cut -d: -f1 | sort -u); do sed -i'' 's,'"\"$x/"',"github.com/tiborvass/gomod/internal/'${x%internal}',g' $f; done; done

go list -json . | jq -r .Deps[] | grep github.com/tiborvass/gomod | sed 's,^github.com/tiborvass/gomod/internal/,,g' | sort -u > /tmp/1
( cd internal && find . -type d -not -name internal -not -name cmd -not -name go | cut -d/ -f2- | sort -u > /tmp/2 )
rm -rf $(diff -u /tmp/{1,2} | grep '^+i' | cut -d+ -f2)
