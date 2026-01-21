package openapi

import (
	"net/http"
	"strings"

	"github.com/aizacoders/openapigo/openapi/infer"

	"github.com/getkin/kin-openapi/openapi3"
)

// BuildSpec builds an OpenAPI document from captured routes and config.
func BuildSpec(routes []RouteMeta, cfg Config) *openapi3.T {
	doc := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:   cfg.Title,
			Version: cfg.Version,
		},
		Paths: openapi3.NewPaths(),
		Components: &openapi3.Components{
			Schemas:         map[string]*openapi3.SchemaRef{},
			SecuritySchemes: openapi3.SecuritySchemes{},
		},
	}

	// Security schemes
	if cfg.SecuritySchemes != nil {
		for k, v := range cfg.SecuritySchemes {
			doc.Components.SecuritySchemes[k] = v
		}
	}

	for _, route := range routes {
		path := infer.NormalizePath(route.Path)
		op := &openapi3.Operation{
			Summary:     firstNonEmpty(route.Summary, route.Path),
			Description: route.Description,
			Responses:   &openapi3.Responses{},
		}

		// Path parameters
		if len(route.PathParams) > 0 {
			for _, pp := range route.PathParams {
				if strings.TrimSpace(pp.Name) == "" {
					continue
				}
				op.AddParameter(&openapi3.Parameter{
					Name:        pp.Name,
					In:          openapi3.ParameterInPath,
					Required:    pp.Required,
					Description: pp.Description,
					Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapiTypeToSchemaType(pp.Type)}},
				})
			}
		} else {
			for _, p := range infer.PathParams(route.Path) {
				op.AddParameter(&openapi3.Parameter{
					Name:     p,
					In:       openapi3.ParameterInPath,
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
				})
			}
		}

		// Query parameters (declared via WithQueryParams)
		if len(route.QueryParams) > 0 {
			addQueryParams(op, route.QueryParams)
		}

		if route.RequestSchema != nil {
			schemaRef := infer.RequestSchema(doc, route.RequestSchema)
			op.RequestBody = &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{Required: true, Content: openapi3.NewContentWithJSONSchemaRef(schemaRef)}}
		}

		if route.ResponseSchema != nil {
			schemaRef := infer.ResponseSchema(doc, route.ResponseSchema)
			op.Responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK"), Content: openapi3.NewContentWithJSONSchemaRef(schemaRef)}})
		} else {
			op.Responses.Set("200", &openapi3.ResponseRef{Value: &openapi3.Response{Description: ptr("OK")}})
		}

		if route.Security != nil {
			op.Security = &openapi3.SecurityRequirements{*route.Security}
		}

		item := &openapi3.PathItem{}
		switch route.Method {
		case http.MethodGet:
			item.Get = op
		case http.MethodPost:
			item.Post = op
		case http.MethodPut:
			item.Put = op
		case http.MethodDelete:
			item.Delete = op
		case http.MethodPatch:
			item.Patch = op
		case http.MethodHead:
			item.Head = op
		case http.MethodOptions:
			item.Options = op
		case http.MethodTrace:
			item.Trace = op
		}

		doc.Paths.Set(path, item)
	}

	return doc
}

func firstNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
