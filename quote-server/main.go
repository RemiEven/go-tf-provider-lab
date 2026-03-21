package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

//go:embed index.html
var htmlContent string

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	t := template.Must(template.New("tmpl").Parse(htmlContent))

	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		quotes, err := getQuotes()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, quotes)
	})

	fmt.Println("will listen to port 5678")

	http.ListenAndServe(":5678", nil)

	return nil
}

type Quote struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

func getQuotes() ([]Quote, error) {
	fileSystem := os.DirFS("..")
	readDirFileSystem, ok := fileSystem.(fs.ReadDirFS)
	if !ok {
		return nil, errors.New("os.DirFS returned a value that doesn't implement fs.ReadDirFS")
	}

	const quotesFolder = "myquotes"

	dirEntries, err := readDirFileSystem.ReadDir(quotesFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read quotes folder: %w", err)
	}

	quotes := make([]Quote, 0, len(dirEntries))
	for _, entry := range dirEntries {
		file, err := readDirFileSystem.Open(quotesFolder + "/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open quote file: %w", err)
		}
		defer file.Close()

		var quote Quote
		if err := json.NewDecoder(file).Decode(&quote); err != nil {
			return nil, fmt.Errorf("failed to decode quote file %s: %w", entry.Name(), err)
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}
