package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Options rappresenta le opzioni raccolte dalla CLI.
type Options struct {
	JSONPath   string
	SpecPath   string
	SchemaName string
	Path       string
	Method     string
}

// Parse legge gli argomenti della CLI e restituisce le Options.
func Parse(args []string) (Options, error) {
	var opts Options

	fs := flag.NewFlagSet("oas-validator", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.SchemaName, "schema", "", "nome dello schema in components.schemas da usare")
	fs.StringVar(&opts.Path, "path", "", "path dell'endpoint REST da cui estrarre lo schema")
	fs.StringVar(&opts.Method, "method", "", "metodo HTTP dell'endpoint REST")
	help := fs.Bool("h", false, "mostra l'uso")
	helpLong := fs.Bool("help", false, "mostra l'uso")

	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if *help || *helpLong {
		return opts, errors.New("aiuto richiesto")
	}

	remaining := fs.Args()
	if len(remaining) != 2 {
		return opts, errors.New("sono richiesti <file.json> e <spec.{yaml|yml|json}>")
	}

	opts.JSONPath = remaining[0]
	opts.SpecPath = remaining[1]
	opts.Method = strings.ToUpper(opts.Method)

	if opts.SchemaName != "" && (opts.Path != "" || opts.Method != "") {
		return opts, errors.New("usare --schema oppure --path/--method, non entrambi")
	}
	if (opts.Path != "" && opts.Method == "") || (opts.Path == "" && opts.Method != "") {
		return opts, errors.New("--path e --method devono essere usati insieme")
	}

	return opts, nil
}

// PrintUsage scrive sull'writer indicato l'uso della CLI.
func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, "Uso:\n oas-validator [--schema NomeSchema | --path /percorso --method POST] <file.json> <spec.{yaml|yml|json}>")
}

// MaybeSuggest stampa suggerimenti su schemi e path disponibili quando la selezione fallisce.
func MaybeSuggest(w io.Writer, doc *openapi3.T, opts Options) {
	if doc == nil {
		return
	}
	var sections []string

	if doc.Components.Schemas != nil && len(doc.Components.Schemas) > 0 {
		names := make([]string, 0, len(doc.Components.Schemas))
		for k := range doc.Components.Schemas {
			names = append(names, k)
		}
		sort.Strings(names)
		if len(names) > 0 {
			sections = append(sections, "Suggerimento: disponibili in components.schemas -> "+strings.Join(names, ", "))
		}
	}

	if opts.Path == "" && opts.Method == "" && doc.Paths != nil {
		paths := doc.Paths.Map()
		rows := make([]string, 0, len(paths))
		for p, item := range paths {
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
				rows = append(rows, fmt.Sprintf("%s [%s]", p, strings.Join(methods, ",")))
			}
		}
		sort.Strings(rows)
		if len(rows) > 0 {
			sections = append(sections, "Oppure usa --path/--method. Paths disponibili:\n- "+strings.Join(rows, "\n- "))
		}
	}

	if len(sections) > 0 {
		fmt.Fprintln(w, strings.Join(sections, "\n"))
	}
}

// PrintOK stampa un messaggio di riepilogo quando la validazione va a buon fine.
func PrintOK(w io.Writer, doc *openapi3.T, schema *openapi3.SchemaRef, opts Options) {
	name := "<schema>"
	if doc != nil && doc.Components.Schemas != nil {
		for k, v := range doc.Components.Schemas {
			if v == schema {
				name = "components.schemas." + k
				break
			}
		}
	}
	if name == "<schema>" {
		name = "requestBody application/json"
	}
	fmt.Fprintf(w, "OK (%s) → %s\n", filepathBase(name), filepathBase(opts.JSONPath))
}

func filepathBase(path string) string {
	if path == "" {
		return path
	}
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}
	if idx := strings.LastIndex(path, "\\"); idx >= 0 {
		return path[idx+1:]
	}
	return path
}
