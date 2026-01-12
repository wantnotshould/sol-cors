# sol-cors

For [sol](https://github.com/wantnotshould/sol)

## Install

```bash
go get github.com/wantnotshould/sol-cors
```

## Quick Start

```go
package main

import (
	"github.com/wantnotshould/sol"
	"github.com/wantnotshould/sol-cors"
)

func main(){
    sl := sol.New()
    sl.Use(cors.Default())

    sl.POST("/", func(c *sol.Context) {
		fmt.Fprintln(c.Writer, "/")
	})

	sl.Run()
}
```