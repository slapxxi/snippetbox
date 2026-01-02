package main

import (
	"html/template"
	"path/filepath"

	"github.com/slapxxi/snippetbox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/index.html",
			"./ui/html/partials/nav.html",
			page,
		}

		templ, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = templ
	}

	return cache, nil
}
