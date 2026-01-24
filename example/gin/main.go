//go:build gin && !typed && !security

package main

import (
	ginadapter "github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi/simple"
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
	sr := simple.NewGinRouter(r, openapiSpec())

	registerSystemRoutes(sr)
	registerUserRoutes(sr)

	ginadapter.Register(r, openAPICfg())
	_ = r.Engine.Run(":8080")
}
