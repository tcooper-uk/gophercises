package template

import (
	"adventure/internal/story"
	"html/template"
	"io"
)

func ServeStoryTemplate(writer io.Writer, story *story.Story) error {

	tmpl := template.Must(template.ParseFiles("web/template/layout.html"))
	err := tmpl.Execute(writer, story)

	if err != nil {
		return err
	}

	return nil
}
