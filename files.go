package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// fileExists verifica se un percorso esiste ed è accessibile.
func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// mustArgFiles controlla che i due percorsi ai file esistano; in caso contrario restituisce un errore.
// Usato per emettere exit code 2 quando i file sono assenti.
func mustArgFiles(jsonPath, specPath string) error {
	if jsonPath == "" || specPath == "" {
		return errors.New("file assenti o argomenti mancanti")
	}
	if !fileExists(jsonPath) || !fileExists(specPath) {
		return fmt.Errorf("file assenti: %s oppure %s", jsonPath, specPath)
	}
	return nil
}

// loadJSON legge e deserializza un file JSON generico in una variabile di tipo any.
// Verifica anche la correttezza sintattica del JSON.
func loadJSON(path string) (any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("JSON non valido: %w", err)
	}
	return v, nil
}
