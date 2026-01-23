# Fiber example (OpenAPIGO)

Fiber example uses the same config-first style with `openapi/simple`.

## Quick start

Install Fiber if you don't have it:

```bash
go get github.com/gofiber/fiber/v2@latest
```

Run the example:

```bash
go run ./example/fiber
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./example/fiber
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
    fiberlib "github.com/gofiber/fiber/v2"
    fiberadapter "github.com/aizacoders/openapigo/adapters/fiber"
    "github.com/aizacoders/openapigo/openapi"
    "github.com/aizacoders/openapigo/openapi/simple"
)
```

2) Create Fiber app and wrap with adapter

```go
app := fiberlib.New()
adapter := fiberadapter.NewFromApp(app)
```

3) Build Spec with `simple.NewSpec()`

```go
b := simple.NewSpec()
b.GroupTags("/", []string{"Users"}, func(s *simple.SpecBuilder) {
    s.GET("/users").Res([]User{}).OK()
    s.POST("/users").Req(CreateUser{}).Res(User{}).Created()
})
```

4) Create simple wrapper and register handlers

```go
sr := simple.NewFiber(adapter, b.Spec())
users := sr.Group("", fiberadapter.WithTags("Users"))
users.GET("/users", func(c *fiberlib.Ctx) error {
    return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
})
```

5) Mount OpenAPI and run

```go
adapter.Register(adapter, openapi.Config{Title: "User API", Version: "1.0.0"})
adapter.App.Listen(":8080")
```

6) Notes

- `NewFromApp` allows you to configure middleware and settings on the Fiber app before wrapping it with the adapter.
- Use `MultipartUpload` in the Spec builder to expose file upload in Swagger UI.
