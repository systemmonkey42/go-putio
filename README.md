# putio  [![Build Status](https://travis-ci.org/igungor/go-putio.svg?branch=master)](https://travis-ci.org/igungor/go-putio)

[](https://godoc.org/github.com/igungor/go-putio/putio)

putio is a Go client library for accessing the Put.io v2 API.

## status

API is not stable yet.

## install

```sh
go get github.com/igungor/go-putio"
```

## usage

```go
package main

import "github.com/igungor/go-putio/putio"

func main() {
    oauthClient := putio.NewAuthHelper("YOUR-TOKEN-HERE")
    client := putio.NewClient(oauthClient)
    resp, _ := client.Get(0)
    println(resp.Parent.Filename)
}
```
