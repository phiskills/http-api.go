# Phi Skills Generic HTTP API Server for Go

| **Homepage** | [https://phiskills.com][0]        |
| ------------ | --------------------------------- | 
| **GitHub**   | [https://github.com/phiskills][1] |

## Overview

This project contains the Go module to create a generic **HTTP API Server**.  

## Installation

```bash
go get github.com/phiskills/http-api.go
```

## Creating the server

```go
package main
import "github.com/phiskills/http-api.go"

api := http.New('My API')
api.Register("/", http.Router{
	Get: func(ctx *http.Context) { ... },
})
api.Start()
```
For more details, see [Package http][10].

[0]: https://phiskills.com
[1]: https://github.com/phiskills
[10]: https://golang.org/pkg/net/http/
