package story

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

func serveStoryTemplate(writer io.Writer, story *Story, layout string) error {

	if !strings.HasSuffix(layout, ".html") {
		layout = fmt.Sprintf("%s.html", layout)
	}

	path := fmt.Sprintf("web/template/%s", layout)

	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(writer, story)

	if err != nil {
		return err
	}

	return nil
}
