package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/the-zucc/typekit"
	"github.com/the-zucc/typekit/example/config"
)

type CustomHttpServer struct {
	serverPortStr string
}

func (s *CustomHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func (s *CustomHttpServer) ServeSync() error {
	fmt.Printf("serving HTTP - %s", s.serverPortStr)
	return http.ListenAndServe(s.serverPortStr, s)
}

// here we register the instance of our server type,
// for use with typekit.Resolve[]()
var _ = typekit.Register(func() (CustomHttpServer, error) {
	return CustomHttpServer{
		serverPortStr: typekit.Resolve[config.MyAppConfig]().ServerPortStr,
	}, nil
})
