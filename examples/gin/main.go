//go:build gin && !typed && !security

package main

import (
	"github.com/aizacoders/openapigo/adapters/ginadapter"
	"github.com/aizacoders/openapigo/openapi/oas"
	ginlib "github.com/gin-gonic/gin"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	engine := ginlib.New()

	// wrap existing engine into adapter router so OpenAPI metadata is captured
	r := ginadapter.NewGinAdapters(engine)
	sr := oas.NewGinRouter(r, openapiSpec())

	registerSystemRoutes(sr)
	registerUserRoutes(sr)

	ginadapter.Register(r, openAPICfg())
	_ = r.Engine.Run(":8080")
}
