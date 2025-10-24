package main

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

// validateAgainst valida un dato generico (mappa/array derivato da JSON) contro uno schema OpenAPI.
// Usa un validator consapevole del documento per risolvere correttamente i $ref.
func validateAgainst(doc *openapi3.T, schema *openapi3.SchemaRef, data any) error {
	validator := openapi3.NewSchemaValidator(schema, doc, "", nil)
	return validator.Validate(context.Background(), data)
}
