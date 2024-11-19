package main

import (
	"github.com/the-zucc/typekit"
	"github.com/the-zucc/typekit/example/server"
)

func main() {
	typekit.Get[server.CustomHttpServer]().ServeSync()
}
