package config

import (
	"os"

	"github.com/the-zucc/typekit"
)

type MyAppConfig struct {
	ServerPortStr string
	Protocol      string
}

// here we register the instance of MyAppConfig
// that is to be used alongside typekit.Resolve[]()
var config = typekit.Register(MyAppConfig{
	ServerPortStr: getenv("SERVER_PORT", ":8080"),
	Protocol:      getenv("PROTOCOL", "http"),
})

func getenv(key string, default_ string) string {
	if val := os.Getenv(key); val == "" {
		return default_
	} else {
		return val
	}
}
