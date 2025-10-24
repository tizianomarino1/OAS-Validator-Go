package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"oas-validator-go/internal/cli"
	"oas-validator-go/internal/files"
	"oas-validator-go/internal/spec"
	"oas-validator-go/internal/validation"
)

const logFilenameTemplate = "oas-validator-log-%s.txt"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) (code int) {
	logPath := fmt.Sprintf(logFilenameTemplate, time.Now().Format("20060102-150405"))
	logFile, err := os.Create(logPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "impossibile creare il file di log:", err)
		return 1
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("Avvio oas-validator con argomenti: %v", args)
	defer func() {
		logger.Printf("Terminazione con codice: %d", code)
	}()

	opts, err := cli.Parse(args)
	if err != nil {
		logger.Printf("Errore durante il parsing delle opzioni: %v", err)
		fmt.Fprintln(os.Stderr, err)
		cli.PrintUsage(os.Stderr)
		printLogHint(os.Stderr, logPath)
		return 1
	}
	logger.Printf("Opzioni ricevute: JSONPath=%s SpecPath=%s SchemaName=%s Path=%s Method=%s", opts.JSONPath, opts.SpecPath, opts.SchemaName, opts.Path, opts.Method)

	if err := files.MustExist(opts.JSONPath, opts.SpecPath); err != nil {
		logger.Printf("File mancanti: %v", err)
		fmt.Fprintln(os.Stderr, err)
		printLogHint(os.Stderr, logPath)
		return 2
	}
	logger.Println("Entrambi i file sono stati trovati")

	data, err := files.LoadJSON(opts.JSONPath)
	if err != nil {
		logger.Printf("Errore caricando il JSON: %v", err)
		fmt.Fprintln(os.Stderr, err)
		printLogHint(os.Stderr, logPath)
		if isMissingFileError(err) {
			return 2
		}
		return 1
	}
	logger.Println("JSON caricato correttamente")

	doc, err := spec.Load(opts.SpecPath)
	if err != nil {
		logger.Printf("Errore caricando la specifica: %v", err)
		fmt.Fprintln(os.Stderr, err)
		printLogHint(os.Stderr, logPath)
		if isMissingFileError(err) {
			return 2
		}
		return 1
	}
	logger.Println("Specifiche OpenAPI caricate correttamente")

	schema, err := spec.PickSchema(doc, opts.SchemaName, opts.Path, opts.Method)
	if err != nil {
		logger.Printf("Errore selezionando lo schema: %v", err)
		fmt.Fprintln(os.Stderr, "Selezione schema:", err)
		cli.MaybeSuggest(os.Stdout, doc, opts)
		printLogHint(os.Stderr, logPath)
		return 1
	}
	logger.Println("Schema selezionato con successo")

	if err := validation.ValidateAgainst(doc, schema, data); err != nil {
		logger.Printf("Validazione fallita: %v", err)
		fmt.Fprintln(os.Stderr, "Validazione FALLITA:")
		fmt.Fprintln(os.Stderr, err)
		printLogHint(os.Stderr, logPath)
		return 1
	}
	logger.Println("Validazione completata con successo")

	cli.PrintOK(os.Stdout, doc, schema, opts)
	printLogHint(os.Stdout, logPath)
	return 0
}

func printLogHint(w *os.File, path string) {
	fmt.Fprintf(w, "Dettagli disponibili nel log: %s\n", path)
}

func isMissingFileError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, os.ErrNotExist) {
		return true
	}
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return errors.Is(pathErr.Err, fs.ErrNotExist)
	}
	return false
}
