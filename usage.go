package main

import (
	"fmt"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

// printUsage stampa una breve guida d'uso.
func printUsage() {
	fmt.Println("Uso:\n oas-validator [--schema NomeSchema | --path /percorso --method POST] <file.json> <spec.{yaml|yml|json}>")
}

// maybeSuggest stampa suggerimenti su schemi e path disponibili quando la selezione fallisce.
func maybeSuggest(doc *openapi3.T, opts Options) {
	listHint := ""
	if doc.Components.Schemas != nil && len(doc.Components.Schemas) > 0 {
		names := make([]string, 0, len(doc.Components.Schemas))
		for k := range doc.Components.Schemas {
			names = append(names, k)
		}
		if len(names) > 0 {
			listHint = "\nSuggerimento: disponibili in components.schemas -> " + stringsJoin(names, ", ")
		}
	}
	if opts.Path == "" && opts.Method == "" && doc.Paths != nil {
		rows := make([]string, 0, len(doc.Paths))
		for p, item := range doc.Paths {
			if item == nil {
				continue
			}
			methods := make([]string, 0, 8)
			if item.Get != nil {
				methods = append(methods, "GET")
			}
			if item.Post != nil {
				methods = append(methods, "POST")
			}
			if item.Put != nil {
				methods = append(methods, "PUT")
			}
			if item.Patch != nil {
				methods = append(methods, "PATCH")
			}
			if item.Delete != nil {
				methods = append(methods, "DELETE")
			}
			if item.Options != nil {
				methods = append(methods, "OPTIONS")
			}
			if item.Head != nil {
				methods = append(methods, "HEAD")
			}
			if item.Trace != nil {
				methods = append(methods, "TRACE")
			}
			if len(methods) > 0 {
				rows = append(rows, fmt.Sprintf("%s [%s]", p, stringsJoin(methods, ",")))
			}
		}
		if len(rows) > 0 {
			listHint += "\nOppure usa --path/--method. Paths disponibili:\n- " + stringsJoin(rows, "\n- ")
		}
	}
	if listHint != "" {
		fmt.Println(listHint)
	}
}

// printOK stampa un messaggio di riepilogo quando la validazione è andata a buon fine.
func printOK(doc *openapi3.T, schema *openapi3.SchemaRef, opts Options) {
	// Prova a ricavare un nome "amichevole" per lo schema utilizzato
	name := "<schema>"
	for k, v := range doc.Components.Schemas {
		if v == schema {
			name = "components.schemas." + k
			break
		}
	}
	if name == "<schema>" {
		name = "requestBody application/json"
	}
	fmt.Printf("OK (%s) → %s\n", filepath.Base(name), filepath.Base(opts.JSONPath))
}

// stringsJoin è un piccolo wrapper per evitare un'import separata di strings in questo file.
func stringsJoin(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	out := elems[0]
	for i := 1; i < len(elems); i++ {
		out += sep + elems[i]
	}
	return out
}
