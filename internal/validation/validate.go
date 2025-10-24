package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ValidateAgainst valida un dato generico contro uno schema OpenAPI.
func ValidateAgainst(doc *openapi3.T, schema *openapi3.SchemaRef, data any) error {
	resolved, err := resolveSchema(doc, schema)
	if err != nil {
		return err
	}
	return resolved.VisitJSON(data, openapi3.MultiErrors())
}

func resolveSchema(doc *openapi3.T, ref *openapi3.SchemaRef) (*openapi3.Schema, error) {
	if ref == nil {
		return nil, errors.New("schema nullo")
	}
	if ref.Value != nil {
		return ref.Value, nil
	}
	if ref.Ref == "" {
		return nil, errors.New("schema senza definizione")
	}
	if doc != nil {
		const prefix = "#/components/schemas/"
		if strings.HasPrefix(ref.Ref, prefix) {
			name := strings.TrimPrefix(ref.Ref, prefix)
			if target, ok := doc.Components.Schemas[name]; ok && target != nil && target.Value != nil {
				return target.Value, nil
			}
		}
	}
	return nil, fmt.Errorf("impossibile risolvere lo schema %q", ref.Ref)
}
