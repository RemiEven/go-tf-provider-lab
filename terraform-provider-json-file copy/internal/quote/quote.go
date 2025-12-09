package quote

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

type Quote struct {
	Message string `json:"message"`
	Author  string `json:"author,omitempty"`
}

func CreateQuoteFile(folder string, quote Quote) (string, error) {
	id := uuid.New().String()

	if err := WriteQuoteFile(folder, id, quotei); err != nil {
		return "", err
	}
	return id, nil
}

func WriteQuoteFile(folder string, id string, quote Quote) error {
	file, err := os.Create(filename(folder, id))
	if err != nil {
		return fmt.Errorf("failed to create/truncate quote file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(quote); err != nil {
		return fmt.Errorf("failed to write quote in file: %w", err)
	}

	return nil
}

func ReadQuote(folder string, id string) (*Quote, error) {
	file, err := os.Open(filename(folder, id))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read quote file: %w", err)
	}
	defer file.Close()

	var q Quote

	if err := json.NewDecoder(file).Decode(&q); err != nil {
		return nil, fmt.Errorf("failed to decode quote file: %q", err)
	}
	return &q, nil
}

func DeleteQuoteFile(folder string, id string) error {
	if err := os.Remove(filename(folder, id)); err != nil {
		return fmt.Errorf("failed to delete quote file: %w", err)
	}
	return nil
}

func filename(folder string, id string) string {
	return path.Join(folder, id+".json")
}
