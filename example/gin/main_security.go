//go:build gin && !typed && security

package main

import (
	"github.com/aizacoders/openapigo/adapters/gin"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	cfg, bearer, apiKey := openAPICfgSecurity()
	registerSecureRoutes(r, bearer, apiKey)

	gin.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
