package openapi

import "net/http"

// JSONRouteSpec is a convenience for common JSON APIs.
// It wires request/response schemas + a primary success status code.
//
// You still can override everything by passing explicit options.
//
// Typical usage from adapters:
//
//	r.POST("/users", h, openapi.JSONRoute(CreateUser{}, struct{}{}, http.StatusCreated)...)
func JSONRoute(reqSchema any, resSchema any, successStatus int) []HandlerOption {
	opts := make([]HandlerOption, 0, 3)
	if reqSchema != nil {
		opts = append(opts, WithRequestSchema(reqSchema))
	}
	if resSchema != nil {
		opts = append(opts, WithResponseSchema(resSchema))
	}
	if successStatus == 0 {
		successStatus = http.StatusOK
	}
	// declare the primary success response; default errors (400/500/401) are handled by builder when WithResponses isn't used.
	opts = append(opts, WithResponses(ResponseSpec{Status: successStatus, Schema: resSchema}))
	return opts
}
