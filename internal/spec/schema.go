package spec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// PickSchema seleziona lo schema da usare per la validazione in base alle opzioni fornite.
func PickSchema(doc *openapi3.T, schemaName, path, method string) (*openapi3.SchemaRef, error) {
	if doc == nil {
		return nil, errors.New("documento OpenAPI nullo")
	}

	if schemaName != "" {
		if schema := getSchemaByName(doc, schemaName); schema != nil {
			return schema, nil
		}
		return nil, fmt.Errorf("schema %q non trovato in components.schemas", schemaName)
	}

	if path != "" || method != "" {
		if path == "" || method == "" {
			return nil, errors.New("specificare sia --path che --method")
		}
		if schema := findSchemaFromPathMethod(doc, path, method); schema != nil {
			return schema, nil
		}
		return nil, fmt.Errorf("nessun requestBody application/json per %s %s", strings.ToUpper(method), path)
	}

	if schema := findUniqueReqBodyJSONSchema(doc); schema != nil {
		return schema, nil
	}

	if doc.Components.Schemas != nil {
		var fallback *openapi3.SchemaRef
		for name, schema := range doc.Components.Schemas {
			if schema == nil {
				continue
			}
			if preferredSchemaName(name) {
				return schema, nil
			}
			if fallback == nil {
				fallback = schema
			}
		}
		if fallback != nil {
			return fallback, nil
		}
	}

	return nil, errors.New("impossibile determinare lo schema da usare")
}

func preferredSchemaName(name string) bool {
	n := strings.ToLower(name)
	return n == "instancedescriptor" || n == "body"
}

func getSchemaByName(doc *openapi3.T, name string) *openapi3.SchemaRef {
	if doc.Components.Schemas == nil {
		return nil
	}
	if s, ok := doc.Components.Schemas[name]; ok && s != nil {
		return s
	}
	lname := strings.ToLower(name)
	for k, v := range doc.Components.Schemas {
		if strings.ToLower(k) == lname {
			return v
		}
	}
	return nil
}

func findUniqueReqBodyJSONSchema(doc *openapi3.T) *openapi3.SchemaRef {
	found := map[*openapi3.SchemaRef]struct{}{}
	var picked *openapi3.SchemaRef

	collect := func(op *openapi3.Operation) {
		if op == nil || op.RequestBody == nil || op.RequestBody.Value == nil {
			return
		}
		if mt, ok := op.RequestBody.Value.Content["application/json"]; ok && mt != nil && mt.Schema != nil {
			if _, seen := found[mt.Schema]; !seen {
				found[mt.Schema] = struct{}{}
				picked = mt.Schema
			}
		}
	}

	if doc.Paths != nil {
		for _, p := range doc.Paths.Map() {
			if p == nil {
				continue
			}
			collect(p.Get)
			collect(p.Post)
			collect(p.Put)
			collect(p.Patch)
			collect(p.Delete)
			collect(p.Options)
			collect(p.Head)
			collect(p.Trace)
		}
	}

	if len(found) == 1 {
		return picked
	}
	return nil
}

func findSchemaFromPathMethod(doc *openapi3.T, path, method string) *openapi3.SchemaRef {
	if doc.Paths == nil {
		return nil
	}
	item := doc.Paths.Value(path)
	if item == nil {
		return nil
	}

	var op *openapi3.Operation
	switch strings.ToUpper(method) {
	case "GET":
		op = item.Get
	case "POST":
		op = item.Post
	case "PUT":
		op = item.Put
	case "PATCH":
		op = item.Patch
	case "DELETE":
		op = item.Delete
	case "OPTIONS":
		op = item.Options
	case "HEAD":
		op = item.Head
	case "TRACE":
		op = item.Trace
	default:
		return nil
	}

	if op == nil || op.RequestBody == nil || op.RequestBody.Value == nil {
		return nil
	}
	if mt, ok := op.RequestBody.Value.Content["application/json"]; ok && mt != nil {
		return mt.Schema
	}
	return nil
}
