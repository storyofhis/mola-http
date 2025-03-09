# MOLA
tree based routing protocol 
example of this code : 

```
package main

import (
	"log"

	"github.com/storyofhis/mola/mola-http"
)

func main() {
	server := mola.NewServer()

	// Define routes
	server.GET("/", func(ctx *mola.Context) {
		ctx.Text(200, "Welcome to Mola HTTP framework!")
	}, mola.LoggerMiddleware)

	server.GET("/user/:name", func(ctx *mola.Context) {
		name := ctx.Params["name"]
		ctx.JSON(200, map[string]string{"message": "Hello, " + name})
	}, mola.LoggerMiddleware)

	server.POST("/user/:name", func(ctx *mola.Context) {
		name := ctx.Params["name"]
		ctx.JSON(200, map[string]string{"message": "Hello, " + name})
	}, mola.LoggerMiddleware)

	// Start the server
	log.Println("Server running on port 8080")
	err := server.Start("8080")
	if err != nil {
		log.Fatal(err)
	}
}

```
