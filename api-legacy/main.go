/*
File: teamknowl/api/main.go
Purpose: Core REST API for TeamKnowl, serving Markdown documentation to AI agents and the UI.
Product/business importance: Provides the primary interface for both human and AI interactions with the synchronized knowledge base.

Copyright (c) 2026 John K Johansen
License: MIT (see LICENSE)
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "teamknowl/api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title TeamKnowl API
// @version 1.0
// @description agent-friendly API for serving Markdown documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name John K Johansen
// @contact.url http://johnkjohansen.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// Config holds the API configuration.
type Config struct {
	DocsDir string
	Port    string
}

// FileInfo represents a Markdown file in the knowledge base.
type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func main() {
	config := Config{
		DocsDir: getEnv("DOCS_DIR", "/docs"),
		Port:    getEnv("PORT", "8080"),
	}

	http.HandleFunc("/v1/list", listFiles(config))
	http.HandleFunc("/v1/context", getContext(config))
	http.HandleFunc("/healthz", healthCheck)

	// Swagger documentation
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Printf("TeamKnowl API starting on port %s, serving from %s", config.Port, config.DocsDir)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}

// listFiles returns a list of all Markdown files in the configured directory.
// @Summary List documentation files
// @Description get all markdown files available in the knowledge base
// @Tags documentation
// @Accept  json
// @Produce  json
// @Success 200 {array} FileInfo
// @Router /v1/list [get]
func listFiles(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var files []FileInfo
		err := filepath.Walk(cfg.DocsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip hidden directories (except .worktrees which git-sync uses)
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != ".worktrees" {
				return filepath.SkipDir
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				relPath, _ := filepath.Rel(cfg.DocsDir, path)
				files = append(files, FileInfo{
					Name: info.Name(),
					Path: relPath,
				})
			}
			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	}
}

// getContext returns the content of a specific Markdown file for AI context.
// @Summary Get documentation context
// @Description get the full markdown content of a specific file
// @Tags documentation
// @Accept  json
// @Produce  plain
// @Param path query string true "Relative path to the markdown file"
// @Success 200 {string} string "Markdown content"
// @Failure 404 {string} string "File not found"
// @Router /v1/context [get]
func getContext(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			http.Error(w, "Missing 'path' parameter", http.StatusBadRequest)
			return
		}

		// Ensure the path is safe and within the docs directory.
		fullPath := filepath.Join(cfg.DocsDir, filePath)
		if !strings.HasPrefix(fullPath, filepath.Clean(cfg.DocsDir)) {
			http.Error(w, "Forbidden path", http.StatusForbidden)
			return
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// In a real implementation, we would inject metadata/frontmatter here.
		w.Header().Set("Content-Type", "text/markdown")
		w.Write(content)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
