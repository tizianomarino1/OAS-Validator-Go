package main

import (
	"fmt"
	"os"
)

// main è l'entrypoint dell'applicazione CLI.
// - legge le opzioni dalla riga di comando (parseFlags)
// - valida la presenza dei file (mustArgFiles)
// - carica JSON e OpenAPI (loadJSON, loadOpenAPI)
// - seleziona lo schema (pickSchema)
// - esegue la validazione (validateAgainst)
// Termina con i codici di uscita richiesti: 0 OK, 1 errore di validazione/uso, 2 file assenti o argomenti mancanti.
func main() {
	opts, err := parseFlags(os.Args[1:])
	if err != nil {
		// Errore di parsing argomenti → usage e exit 2 se mancano argomenti o file
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(2)
	}

	if err := mustArgFiles(opts.JSONPath, opts.SpecPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	data, err := loadJSON(opts.JSONPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	doc, err := loadOpenAPI(opts.SpecPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	schema, err := pickSchema(doc, opts.SchemaName, opts.Path, opts.Method)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Selezione schema:", err)
		maybeSuggest(doc, opts)
		os.Exit(1)
	}

	if err := validateAgainst(doc, schema, data); err != nil {
		fmt.Fprintln(os.Stderr, "Validazione FALLITA:")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	printOK(doc, schema, opts)
	os.Exit(0)
}
