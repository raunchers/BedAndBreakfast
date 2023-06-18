package render

import (
	"bytes"
	"github.com/raunchers/BedAndBreakfast/pkg/config"
	"github.com/raunchers/BedAndBreakfast/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// Variable to hold the template cache
	var tc map[string]*template.Template

	// If use cache is true, read the info from the template cache
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// Rebuild the template cache
		tc, _ = CreateTemplateCache()
	}

	// Get the requested template from the cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get the template cache")
	}

	// buffer to hold bytes
	// execute the buffer directly and then write it out. Allows for better error checking
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	// tells if the error comes from the value that is stored in the map
	_ = t.Execute(buf, td)

	// Render the template that was requested
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// Making an empty map
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from the ./templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	// don't care about the index, use _ to ignore it
	// Page will be populated with whatever we get from the slice of strings, the full path name to the template
	for _, page := range pages {
		// name stores the name of the file that ends in *.page.tmpl
		// ex: about.page.tmpl
		name := filepath.Base(page)

		// Parse the file
		// ts == template set
		// Parse the file stored in page, store that in a template named name
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// look for any layouts that exist in that directory
		// A slice of string with all the layouts
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// Adds the layout to the template set
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		// Adds the final resulting template to the map
		myCache[name] = ts
	}

	return myCache, nil
}
