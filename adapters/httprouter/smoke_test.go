package httprouter

import (
	"testing"

	"github.com/aizacoders/openapigo/openapi"
)

func TestHTTPRouterNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatalf("New() returned nil")
	}
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}
