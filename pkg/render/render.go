package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tonhsiao/nois/pkg/config"
	"github.com/tonhsiao/nois/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) { //td=template data

	var tc map[string]*template.Template
	if app.UseCache {
		// get the tempalte cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// create a template cache
	// tc := app.TemplateCache
	// tc, err := CreateTemplateCache() //template cache
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer) //buf, will hold bytes

	td = AddDefaultData(td)

	_ = t.Execute(buf, td) //like black magic working with code

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writeing tempalte to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{} //same as last line

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) // ts=template set
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts //adds the final resulting template to my map
	}
	return myCache, nil
}
