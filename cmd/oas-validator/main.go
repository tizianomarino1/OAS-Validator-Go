package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"oas-validator-go/internal/cli"
	"oas-validator-go/internal/files"
	"oas-validator-go/internal/spec"
	"oas-validator-go/internal/validation"
)

const (
	logFilenameTemplate = "log-%s.txt"

	exitOK           = 0
	exitGenericError = 1
	exitMissingFile  = 2
)

func main() {
	code := run(os.Args[1:])
	fmt.Println(code + 1)
	os.Exit(code)
}

func run(args []string) (code int) {
	logPath := fmt.Sprintf(logFilenameTemplate, time.Now().Format("20060102-150405"))
	logFile, err := os.Create(logPath)
	if err != nil {
		return exitGenericError
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	opts, err := cli.Parse(args)
	if err != nil {
		logError(logger, err)
		return exitGenericError
	}

	if err := files.MustExist(opts.JSONPath, opts.SpecPath); err != nil {
		logError(logger, err)
		return exitMissingFile
	}

	data, err := files.LoadJSON(opts.JSONPath)
	if err != nil {
		logError(logger, err)
		if isMissingFileError(err) {
			return exitMissingFile
		}
		return exitGenericError
	}

	doc, err := spec.Load(opts.SpecPath)
	if err != nil {
		logError(logger, err)
		if isMissingFileError(err) {
			return exitMissingFile
		}
		return exitGenericError
	}

	schema, err := spec.PickSchema(doc, opts.SchemaName, opts.Path, opts.Method)
	if err != nil {
		logError(logger, err)
		return exitGenericError
	}

	if err := validation.ValidateAgainst(doc, schema, data); err != nil {
		logError(logger, err)
		return exitGenericError
	}

	return exitOK
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

func logError(logger *log.Logger, err error) {
	if logger == nil || err == nil {
		return
	}
	logger.Println(strings.ReplaceAll(err.Error(), "\n", " "))
}
