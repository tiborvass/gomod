#!/bin/bash -e

rm -rf internal
mkdir internal
: ${GOROOT=$(go env GOROOT)}
cp -rf $GOROOT/src/internal/* $GOROOT/src/cmd/internal/* $GOROOT/src/cmd/go/internal/* internal/

for x in '"cmd/go/internal/' '"cmd/internal/' '"internal/'; do for f in $(find . -name '*.go' -exec grep -H "$x" {} \; | cut -d: -f1 | sort -u); do sed -i'' 's,'"$x"',"github.com/tiborvass/gomod/internal/,g' $f; done; done
for f in $(find . -name '*.go' -exec grep -H 'strings.ReplaceAll' {} \; | cut -d: -f1 | sort -u ); do sed -i'' 's,strings.ReplaceAll,mystrings.ReplaceAll,g' $f; goimports -w $f; done

go list -json . | jq -r .Deps[] | grep github.com/tiborvass/gomod | sed 's,^github.com/tiborvass/gomod/,,g' | sort -u > /tmp/1
find internal -mindepth 1 -type d -not -name internal | sort -u > /tmp/2
#diff -u /tmp/{1,2}
rm -rf $(diff -u /tmp/{1,2} | grep '^+i' | cut -d+ -f2)
