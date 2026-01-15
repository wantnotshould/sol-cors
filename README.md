# sol-cors

For [sol](https://github.com/wantnotshould/sol)

> [!WARNING]
> This project is ready for basic testing.
>
> If you wish to update or extend it, please do so in a fork.

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
