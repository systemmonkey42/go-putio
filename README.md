# putio  [![Build Status](https://travis-ci.org/igungor/go-putio.svg?branch=master)](https://travis-ci.org/igungor/go-putio)

putio is a Go client library for accessing the Put.io v2 API.

## status

API is not stable yet.

## documentation

the documentation is available on [godoc](http://godoc.org/github.com/igungor/go-putio/putio).

## install

```sh
go get github.com/igungor/go-putio/putio"
```

## usage

```go
package main

import (
        "fmt"
        "log"

        "golang.org/x/oauth2"
        "github.com/igungor/go-putio/putio"
)

func main() {
    tokenSource := oauth2.StaticTokenSource( &oauth2.Token{AccessToken: "<YOUR-TOKEN-HERE>"})
    oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

    client := putio.NewClient(oauthClient)
    root, err := client.Get(0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(root.Filename)
}
```
