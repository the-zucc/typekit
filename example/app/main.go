package main

import (
	"fmt"

	"github.com/the-zucc/typekit"
	"github.com/the-zucc/typekit/example/server"
)

// this app is a simple web server. We define a custom class,
// of which we manage the single instance using typekit.
func main() {
	server := typekit.Resolve[server.CustomHttpServer]()
	fmt.Printf(
		"error running server - %s",
		// this retrieves the typ
		server.ServeSync(),
	)
}
