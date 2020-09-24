# tgin

[![GoDoc](https://godoc.org/github.com/blacktear23/tgin?status.svg)](https://pkg.go.dev/github.com/blacktear23/tgin?tab=doc)

tgin is a light weight web framework

`tgin` API is most like `gin` web framework. But it is implements by basic go http library. If you focus on small compiled binary size and not very relay on all `gin` features, you can use it instead of `gin`.

# Binary Size

Gin sample code:

```
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		msg := "Hello World"
		value, have := c.GetQuery("name")
		if have {
			msg = "Hello " + value
		}
		c.IndentedJSON(200, gin.H{"Message": msg})
	})
	r.Run("0.0.0.0:8888")
}
```

tgin sample code:

```
package main

import (
	gin "github.com/blacktear23/tgin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		msg := "Hello World"
		value, have := c.GetQuery("name")
		if have {
			msg = "Hello " + value
		}
		c.IndentedJSON(200, gin.H{"Message": msg})
	})
	r.Run("0.0.0.0:8888")
}
```

Compile command:

```
go build -trimpath -ldflags "-s -w"
```

Compiled binary size:

| Package | Binary Size (byte) | Reduce Ratio | go version |
| ------- | ------------------ | ------------ | ---------- |
| gin     | 11863092 Byte      | 0%           | 1.14       |
| tgin    | 6076500 Byte       | 48%          | 1.14       |
| gin     | 11819212 Byte      | 0%           | 1.15       |
| tgin    | 5186060 Byte       | 56%          | 1.15       |
