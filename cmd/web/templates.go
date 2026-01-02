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

		templ, err := template.ParseFiles("./ui/html/index.html")
		if err != nil {
			return nil, err
		}

		templ, err = templ.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		templ, err = templ.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = templ
	}

	return cache, nil
}
