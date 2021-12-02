#!/usr/bin/env bash
export PATH=/usr/local/opt/go/libexec/bin/:$PATH
cd /Users/dmowcomber/Documents/bitbar-plugins/stock-go
GO111MODULE=on go run -mod=vendor stock.go
