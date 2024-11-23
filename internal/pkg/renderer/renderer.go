package renderer

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"todo/internal/util"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
}

// NewRenderer loads templates into a map
func NewRenderer(baseDir string) *TemplateRenderer {
	templates := make(map[string]*template.Template)

	// Walk through all .gohtml files in the base directory
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".gohtml") {
			// Create a normalized name with forward slashes
			relPath, _ := filepath.Rel(baseDir, path)
			name := filepath.ToSlash(relPath) // Convert to forward slashes

			// Parse the template file with base.gohtml included for layout support
			templates[name], err = template.New(name).Funcs(template.FuncMap{
				"statusColor": util.StatusColor, // Register the helper function
			}).ParseFiles(filepath.Join(baseDir, "base.gohtml"), path)
			if err != nil {
				return err
			}
			log.Printf("Template loaded: %s", name) // Log each template
		}
		return nil
	})

	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	return &TemplateRenderer{templates: templates}
}

// Render renders a template using the map
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Get the template from the map
	tmpl, ok := t.templates[name]
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template not found: "+name)
	}

	// Ensure data is a map and pass it to the template
	if data == nil {
		data = map[string]interface{}{}
	}
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse // Example: Add reverse routing support

		// Inject dark mode from context
		user, ok := c.Get("user").(map[string]interface{})
		if ok && user["darkMode"] != nil {
			viewContext["darkMode"] = user["darkMode"].(bool)
		}
	}

	// Render the template into the writer
	return tmpl.ExecuteTemplate(w, "base.gohtml", data)
}
