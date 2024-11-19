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

func (s *CustomHttpServer) ServeSync() {
	fmt.Printf("serving HTTP - %s", s.serverPortStr)
	http.ListenAndServe(s.serverPortStr, s)
}

var appConfig = typekit.Get[config.MyAppConfig]()

var Server = typekit.Register(CustomHttpServer{
	serverPortStr: appConfig.ServerPortStr,
})
