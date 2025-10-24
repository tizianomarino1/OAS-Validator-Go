package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// MustExist controlla che i due percorsi indicati esistano e siano accessibili.
func MustExist(jsonPath, specPath string) error {
	if jsonPath == "" || specPath == "" {
		return errors.New("file assenti o argomenti mancanti")
	}
	if !exists(jsonPath) || !exists(specPath) {
		return fmt.Errorf("file assenti: %s oppure %s", jsonPath, specPath)
	}
	return nil
}

// LoadJSON legge e deserializza un file JSON generico.
func LoadJSON(path string) (any, error) {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
