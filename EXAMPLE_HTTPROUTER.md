# net/http (default) router example (OpenAPIGO)

The "default" router in this repo is `openapi.Router` (a lightweight net/http-based mux).

## Quick start

Run the example:

```bash
go run ./examples/httprouter
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/httprouter
```

Open Swagger UI:
- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:
- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

This section shows how to wire the default HTTP router with OpenAPIGO in your project.

1) Imports

```go
import (
    "net/http"

    "github.com/aizacoders/openapigo/adapters/httprouter"
    "github.com/aizacoders/openapigo/openapi"
    "github.com/aizacoders/openapigo/openapi/oas"
)
```

2) Create the router and build your Spec

```go
// Option A: create a ServeMux and let the adapter mount the router for you
mux := http.NewServeMux()
base := httprouter.New(mux)

// Option B: create router and mount manually
// base := httprouter.New()
// mux := http.NewServeMux()
// mux.Handle("/", base)

b := simple.NewSpec()
b.GroupTags("/", []string{"Users"}, func(s *simple.SpecBuilder) {
    s.GET("/users").Res([]User{}).OK()
    s.POST("/users").Req(CreateUser{}).Res(User{}).Created()
    // multipart upload examples
    s.POST("/users/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()
})
spec := b.Spec()
```

3) Create the `simple` wrapper that injects spec defaults

```go
r := simple.New(base, spec)
```

4) Register handlers using the plain HTTP handler signature

```go
r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
    openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
})

r.POST("/users", func(w http.ResponseWriter, req *http.Request) {
    var in CreateUser
    if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
        openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
        return
    }
    w.WriteHeader(http.StatusCreated)
})
```

5) Mount OpenAPI JSON + Swagger UI on the router and run

```go
httprouter.Register(base, openapi.Config{Title: "User API", Version: "1.0.0"})
_ = http.ListenAndServe(":8080", mux) // serve the ServeMux that has the router mounted
```

6) Security (optional)

- Define schemes in `openapi.Config.SecuritySchemes` and attach per-route via `RouteBuilder.Security` in the Spec builder.

7) Multipart uploads

- Use `MultipartUpload(...)` in the spec to declare a `multipart/form-data` body with a file field; the Swagger UI will render a file chooser and corresponding fields.

---

## Notes

- Examples follow the pattern: build a router/engine, (wrap with adapter when applicable), build spec via `simple.NewSpec()` and then use `simple.New*` wrappers.
- Prefer `simple.New` for net/http example to keep handler signatures standard and clean.
