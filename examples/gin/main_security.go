//go:build gin && security && !typed

package main

import (
	"github.com/aizacoders/openapigo/adapters/ginadapter"
	"github.com/aizacoders/openapigo/openapi/oas"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := ginadapter.New()

	cfg, bearer, apiKey := openAPICfgSecurity()

	sr := oas.NewGinRouter(r, oas.Spec{})
	registerSecureRoutes(sr, bearer, apiKey)

	ginadapter.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
