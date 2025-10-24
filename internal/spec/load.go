package spec

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// Load legge e valida un documento OpenAPI da file.
func Load(path string) (*openapi3.T, error) {
	loader := &openapi3.Loader{Context: context.Background(), IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("caricamento OpenAPI: %w", err)
	}
	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("OpenAPI non valida: %w", err)
	}
	return doc, nil
}
