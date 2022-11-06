package template

import (
	"adventure/internal/story"
	"fmt"
	"html/template"
	"io"
	"strings"
)

func ServeStoryTemplate(writer io.Writer, story *story.Story, layout string) error {

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
