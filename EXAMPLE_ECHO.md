# Echo example (OpenAPIGO)

This example shows the **config-first** (SpringBoot-like) way to use OpenAPIGO with Echo.

## Quick start

Install Echo if you don't have it:

```bash
go get github.com/labstack/echo/v4@latest
```

Run the example:

```bash
go run ./examples/echo
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/echo
```

Open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:

- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

1) Imports

```go
import (
    echolib "github.com/labstack/echo/v4"
    echoadapter "github.com/aizacoders/openapigo/adapters/echo"
    "github.com/aizacoders/openapigo/openapi"
    "github.com/aizacoders/openapigo/openapi/oas"
)
```

2) Create Echo instance and wrap with adapter

```go
base := echolib.New()
adapter := echoadapter.NewEchoAdapters(base)
```

3) Build Spec with `simple.NewSpec()` (group routes, define Req/Res and multipart)

```go
b := simple.NewSpec()
b.GroupTags("/", []string{"Users"}, func(s *simple.SpecBuilder) {
    s.GET("/users").Res([]User{}).OK()
    s.POST("/users").Req(CreateUser{}).Res(User{}).Created()
})
```

4) Create the simple wrapper and register handlers

```go
sr := simple.NewEcho(adapter, b.Spec())
users := sr.Group("", echoadapter.WithTags("Users"))
users.GET("/users", func(c echolib.Context) error {
    return echoadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
})
```

5) Mount OpenAPI and run

```go
adapter.Register(adapter, openapi.Config{Title: "User API", Version: "1.0.0"})
adapter.Echo.Start(":8080")
```

6) Notes

- `NewEchoAdapters` lets you create middleware and configure the Echo instance before wrapping it with the adapter.
- Use `MultipartUpload` in the Spec builder to expose file upload inputs in Swagger UI.

### Note about core router

The OpenAPIGO core router is a lightweight net/http-backed mux. Adapter packages (including Echo) integrate with this core behavior and continue to work as before. If you use the `httprouter` adapter you can optionally mount the router automatically onto a `*http.ServeMux` by calling `httprouter.New(mux)`.

