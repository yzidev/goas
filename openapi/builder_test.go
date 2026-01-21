package openapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

type tUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type tCreateUser struct {
	Name string `json:"name"`
}

func TestRegisterAndSpec(t *testing.T) {
	r := NewRouter()

	r.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		_ = req
		w.WriteHeader(http.StatusOK)
	}, WithResponseSchema(tUser{}))

	jwt := openapi3.NewSecurityRequirement().Authenticate("jwt")
	r.POST("/users", func(w http.ResponseWriter, req *http.Request) {
		var in tCreateUser
		_ = Bind(req, &in)
		w.WriteHeader(http.StatusCreated)
	}, WithRequestSchema(tCreateUser{}), WithSecurity(&jwt))

	Register(r, Config{
		Title:   "Test",
		Version: "0.0.1",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"jwt": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var doc openapi3.T
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if doc.Info == nil || doc.Info.Title != "Test" {
		t.Fatalf("unexpected info: %+v", doc.Info)
	}

	p := doc.Paths.Find("/users/{id}")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET operation for /users/{id}")
	}
	if len(p.Get.Parameters) == 0 {
		t.Fatalf("expected inferred path parameter")
	}

	p2 := doc.Paths.Find("/users")
	if p2 == nil || p2.Post == nil {
		t.Fatalf("expected POST operation for /users")
	}
	if p2.Post.Security == nil || len(*p2.Post.Security) == 0 {
		t.Fatalf("expected security requirement")
	}
}

func TestPathValue(t *testing.T) {
	r := NewRouter()
	r.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		if got := PathValue(req, "id"); got != "123" {
			t.Fatalf("expected path id 123, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestTypedHandlersAutoSchema(t *testing.T) {
	type Req struct {
		Name string `json:"name"`
	}
	type Res struct {
		ID string `json:"id"`
	}

	r := NewRouter()

	POSTT[Req, Res](r, "/typed", func(w http.ResponseWriter, r *http.Request, req Req) (Res, int, error) {
		_ = w
		_ = r
		return Res{ID: req.Name}, http.StatusCreated, nil
	})

	Register(r, Config{Title: "Typed", Version: "1"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	r.ServeHTTP(rec, req)

	var doc openapi3.T
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	p := doc.Paths.Find("/typed")
	if p == nil || p.Post == nil {
		t.Fatalf("expected POST operation for /typed")
	}
	if p.Post.RequestBody == nil || p.Post.RequestBody.Value == nil {
		t.Fatalf("expected requestBody")
	}
	if len(p.Post.RequestBody.Value.Content) == 0 {
		t.Fatalf("expected requestBody content schema")
	}
	if p.Post.Responses == nil {
		t.Fatalf("expected responses")
	}
	resp := p.Post.Responses.Value("200")
	if resp == nil || resp.Value == nil {
		t.Fatalf("expected 200 response")
	}
	if len(resp.Value.Content) == 0 {
		t.Fatalf("expected response content schema")
	}
}

func TestPathParamValueInt(t *testing.T) {
	r := NewRouter()
	r.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		p, err := PathParamValue[int](req, "id")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !p.OK || p.Value != 123 {
			t.Fatalf("expected 123, got %+v", p)
		}
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestQueryParamsInSpec(t *testing.T) {
	r := NewRouter()
	r.GET("/search", func(w http.ResponseWriter, req *http.Request) {
		_, _, _ = QueryValue[int](req, "limit")
		w.WriteHeader(http.StatusOK)
	}, WithQueryParams(
		QueryParam{Name: "q", Type: ParamString, Required: true},
		QueryParam{Name: "limit", Type: ParamInteger, Required: false},
	))

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/search")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /search")
	}

	foundQ := false
	foundLimit := false
	for _, pr := range p.Get.Parameters {
		if pr.Value == nil {
			continue
		}
		if pr.Value.In != openapi3.ParameterInQuery {
			continue
		}
		switch pr.Value.Name {
		case "q":
			foundQ = true
			if !pr.Value.Required {
				t.Fatalf("q should be required")
			}
		case "limit":
			foundLimit = true
		}
	}

	if !foundQ || !foundLimit {
		t.Fatalf("missing query params in spec: q=%v limit=%v", foundQ, foundLimit)
	}
}

func TestMultipleMethodsSamePath(t *testing.T) {
	r := NewRouter()

	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.POST("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users")
	if p == nil {
		t.Fatalf("expected /users path")
	}
	if p.Get == nil {
		t.Fatalf("expected GET /users")
	}
	if p.Post == nil {
		t.Fatalf("expected POST /users")
	}
}

func TestOperationTags(t *testing.T) {
	r := NewRouter()
	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, WithTags("Users"))

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /users")
	}
	if len(p.Get.Tags) != 1 || p.Get.Tags[0] != "Users" {
		t.Fatalf("expected tag Users, got %#v", p.Get.Tags)
	}
}

func TestTopLevelTagsFromConfig(t *testing.T) {
	r := NewRouter()
	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, WithTags("Users"))

	doc := BuildSpec(r.Routes(), Config{
		Title:   "T",
		Version: "1",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	if len(doc.Tags) != 1 || doc.Tags[0].Name != "Users" {
		t.Fatalf("expected top-level tags to contain Users, got %#v", doc.Tags)
	}
}

func TestResponsesMultipleStatusCodes(t *testing.T) {
	type Err struct {
		Error string `json:"error"`
	}
	type User struct {
		ID string `json:"id"`
	}

	r := NewRouter()
	r.POST("/users", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}, WithResponses(
		ResponseSpec{Status: http.StatusCreated, Schema: User{}, Description: "Created"},
		ResponseSpec{Status: http.StatusBadRequest, Schema: Err{}, Description: "Bad Request"},
		ResponseSpec{Status: http.StatusInternalServerError, Schema: Err{}, Description: "Internal Server Error"},
	))

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users")
	if p == nil || p.Post == nil {
		t.Fatalf("expected POST /users")
	}
	if p.Post.Responses.Value("201") == nil {
		t.Fatalf("expected 201 response")
	}
	if p.Post.Responses.Value("400") == nil {
		t.Fatalf("expected 400 response")
	}
	if p.Post.Responses.Value("500") == nil {
		t.Fatalf("expected 500 response")
	}
}
