package main

import (
	"fmt"
	"os"

	"oas-validator-go/internal/cli"
	"oas-validator-go/internal/files"
	"oas-validator-go/internal/spec"
	"oas-validator-go/internal/validation"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cli.PrintUsage(os.Stderr)
		os.Exit(2)
	}

	if err := files.MustExist(opts.JSONPath, opts.SpecPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	data, err := files.LoadJSON(opts.JSONPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	doc, err := spec.Load(opts.SpecPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	schema, err := spec.PickSchema(doc, opts.SchemaName, opts.Path, opts.Method)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Selezione schema:", err)
		cli.MaybeSuggest(os.Stdout, doc, opts)
		os.Exit(1)
	}

	if err := validation.ValidateAgainst(doc, schema, data); err != nil {
		fmt.Fprintln(os.Stderr, "Validazione FALLITA:")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.PrintOK(os.Stdout, doc, schema, opts)
	os.Exit(0)
}
