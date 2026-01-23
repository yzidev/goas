# OpenAPIGO

[![CI](https://github.com/aizacoders/openapigo/actions/workflows/ci.yml/badge.svg)](https://github.com/aizacoders/openapigo/actions/workflows/ci.yml)

Auto-generate **OpenAPI 3.x** from your Go route registrations.

The goal is to keep your routing code **clean** (plain `GET/POST/PUT/PATCH/DELETE`) while still producing a good OpenAPI spec + Swagger UI.

---

## What you get

- `GET /openapi.json` (generated OpenAPI document)
- Swagger UI mounted at:
  - `http://localhost:8080/swagger-ui/index.html#/`
  - `/swagger` is kept as a legacy redirect

---

## Key concepts

### 1) Base router (net/http + chi)

Use the built-in router:

- `openapi.NewRouter()` → returns an `http.Handler`
- register routes with `GET/POST/PUT/PATCH/DELETE`

### 2) Config-first spec (SpringBoot-like)

Go handlers don’t expose schema information automatically.
So OpenAPIGO uses a **config-first** approach:

- put route schemas/tags/security/query/header params in one place using `openapi/simple`
- keep your handlers clean and readable

### 3) Multipart upload support

Use `MultipartUpload(...)` to get `multipart/form-data` request bodies and a file upload field in Swagger UI.

---

## Installation

```bash
go get github.com/aizacoders/openapigo@latest
```

---

## Minimal example (net/http)

```go
package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := openapi.NewRouter()

	// 1) Define spec (grouped, clean)
	b := simple.NewSpec()
	b.GroupTags("", []string{"Users"}, func(s *simple.SpecBuilder) {
		s.GET("/users").Res([]User{}).OK()
	})

	// 2) Mount routes (plain net/http handlers)
	r := simple.New(base, b.Spec())
	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	// 3) Register OpenAPI + Swagger UI
	openapi.Register(base, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = http.ListenAndServe(":8080", r)
}
```

---

## Multipart upload example

In your spec:

```go
s.POST("/users/upload").MultipartUpload(
	"file",
	openapi.MultipartField{Name: "note", Type: openapi.ParamString},
).Res(map[string]string{}).OK()
```

In Swagger UI this will show:
- `file` as file chooser
- `note` as text input
- requestBody content type: `multipart/form-data`

---

## Security

You can provide security schemes via `openapi.Config.SecuritySchemes` and attach requirements per-route.
Examples include two schemes:

- **Bearer** JWT (`Authorization: Bearer <token>`)
- **API key** (`X-API-Key: <key>`)

---

## Examples (recommended)

Run examples and open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

### Default (net/http)

- Docs: [`EXAMPLE_HTTPROUTER.md`](./EXAMPLE_HTTPROUTER.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Gin

- Docs: [`EXAMPLE_GIN.md`](./EXAMPLE_GIN.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Echo

- Docs: [`EXAMPLE_ECHO.md`](./EXAMPLE_ECHO.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Fiber

- Docs: [`EXAMPLE_FIBER.md`](./EXAMPLE_FIBER.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

---

## Current support (today)

OpenAPIGO is currently focused on **4 frameworks/router setups**:

1. **net/http (built-in `openapi.Router` based on chi)**
2. **Gin**
3. **Echo**
4. **Fiber**

Notes:
- Other frameworks may be added later, but the repo intentionally stays small and dependency-light.
- Adapters are provided as packages under `adapters/*` so you can use them when needed.
  They are compiled by default and no special build tags are required to use them.
  If you prefer to keep adapter dependencies optional for your project, consider
  shipping adapters as separate modules (e.g. `github.com/aizacoders/openapigo-adapter-gin`) so downstream projects opt-in.

---

## Roadmap / future updates

The direction going forward:

- **Keep the public API simple**:
  - common HTTP methods only: `GET/POST/PUT/PATCH/DELETE`
  - grouping via `Group(...)`
  - OpenAPI metadata via config-first spec (`openapi/simple`)

- **Improve schema inference gradually**:
  - better tag support (`omitempty`, pointer handling)
  - better nested struct handling
  - better multipart documentation

- **Better DX in Swagger UI**:
  - theming improvements
  - cleaner auth UX
  - consistent error schemas

- **Adapter expansion (optional)**:
  - If more frameworks are added, they will follow the same pattern:
    - keep handlers/framework usage idiomatic
    - keep OpenAPIGO integration minimal
    - keep core library independent of adapter dependencies

### Update policy / compatibility

- The project is evolving quickly.
- We aim to keep the **core API stable** (`openapi.Router`, `openapi.Register`, and `openapi/simple`).
- Adapter APIs may change as we simplify integration and keep parity across frameworks.

### Framework support timeline

For now OpenAPIGO only ships examples + adapters for:
- `net/http` (built-in router)
- Gin
- Echo
- Fiber

Additional frameworks are considered **future work** (optional adapters behind build tags).

### How to add another framework (adapter concept)

If you want to support another framework, the recommended approach is:

- Create a new adapter package under `adapters/<framework>`.
- Guard it with a build tag (so the dependency stays optional).
- The adapter should expose a router wrapper similar to the existing ones:
  - register `GET/POST/PUT/PATCH/DELETE`
  - keep grouping if the framework supports groups
  - call `openapi.Router.Handle(...)` / attach `HandlerOption`s in the same way.

For a starting point, check:
- `adapters/gin`
- `adapters/echo`
- `adapters/fiber`

---

## Adapters (how to use with frameworks)

OpenAPIGO provides lightweight adapters for multiple frameworks so you can keep your
handler code clean while still generating OpenAPI and mounting Swagger UI.

Pattern (recommended):

1. Create your framework engine/app (e.g., `gin`, `echo`, `fiber`).
2. Wrap it with the adapter `NewFrom*` helper (so the adapter captures route metadata).
3. Create the `simple` wrapper using the adapter and your `Spec`.
4. Register OpenAPI via the adapter `Register` helper and run the engine/app.

Examples:

- Gin

```go
import (
    ginlib "github.com/gin-gonic/gin"
    ginadapter "github.com/aizacoders/openapigo/adapters/gin"
    "github.com/aizacoders/openapigo/openapi/simple"
)

engine := ginlib.New()
adapter := ginadapter.NewFromEngine(engine)
sr := simple.NewGin(adapter, mySpec)
// register routes on sr ...
ginadapter.Register(adapter, openapi.Config{Title: "My API", Version: "0.1.0"})
adapter.Engine.Run(":8080")
```

- Echo

```go
import (
    echolib "github.com/labstack/echo/v4"
    echoadapter "github.com/aizacoders/openapigo/adapters/echo"
    "github.com/aizacoders/openapigo/openapi/simple"
)

base := echolib.New()
adapter := echoadapter.NewFromEcho(base)
sr := simple.NewEcho(adapter, mySpec)
// register routes on sr ...
echoadapter.Register(adapter, openapi.Config{Title: "My API", Version: "0.1.0"})
adapter.Echo.Start(":8080")
```

- Fiber

```go
import (
    fiberlib "github.com/gofiber/fiber/v2"
    fiberadapter "github.com/aizacoders/openapigo/adapters/fiber"
    "github.com/aizacoders/openapigo/openapi/simple"
)

app := fiberlib.New()
adapter := fiberadapter.NewFromApp(app)
sr := simple.NewFiber(adapter, mySpec)
// register routes on sr ...
fiberadapter.Register(adapter, openapi.Config{Title: "My API", Version: "0.1.0"})
adapter.App.Listen(":8080")
```

Notes:
- The `NewFrom*` helpers let you keep your preferred engine/app initialization (e.g., `gin.Default()`), while still enabling OpenAPIGO to capture route metadata.
- If you previously built with `-tags`, adapters are now compiled by default — no need to use build tags to get adapter implementations.

---

## License

MIT. See [`LICENSE`](./LICENSE).
