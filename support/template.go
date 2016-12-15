package support

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

// Parse a template and render it out
func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}

// Render parses & render a template
func Render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// StaticHandler handles serving static files
func StaticHandler(staticURL string, staticRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		staticFile := req.URL.Path[len(staticURL):]
		if len(staticFile) != 0 {
			f, err := http.Dir(staticRoot).Open(staticFile)
			if err == nil {
				content := io.ReadSeeker(f)
				http.ServeContent(w, req, staticFile, time.Now(), content)
				return
			}
		}
		http.NotFound(w, req)
	}
}
