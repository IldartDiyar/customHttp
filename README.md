

# Custom HTTP Server

This package provides a simple custom HTTP server implementation in Go. It allows you to create and configure a basic HTTP server with handling for incoming requests.

## How It Works

### Server Initialization

1. Import the `customHttp` package.
```
$ go get github.com/IldartDiyar/customHttp
```
2. Create a handler that satisfies the `Handler` interface.
```go
ServeHTTP(http.ResponseWriter, *http.Request)
```
3. Initialize a server with a specified address and the custom handler.
4. Start the server using the `ListenAndServe` function.

```go
package main

import (
	"log"
	"net/http"

	"github.com/IldartDiyar/customHttp"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", h)
	if err := customHttp.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func h(c *gin.Context) {
	c.String(http.StatusOK, "Salem, alem!")
}
```
