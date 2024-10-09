#!/usr/bin/env sh

export CUR="github.com/Light2Dark/go-starter"
export NEW="github.com/Light2Dark/splitpay"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;