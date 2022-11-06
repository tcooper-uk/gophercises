package main

import (
	"adventure/internal/story"
	"adventure/internal/template"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	inputFile := flag.String("input", "gophers.json", "The JSON input file containing the stories. Default is gophers.json")
	templateName := flag.String("template", "layout.html", "The name of the template to render. Default is layout.html")

	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Unable to read %s\n", *inputFile)
		os.Exit(1)
	}
	defer file.Close()

	json, err := ioutil.ReadAll(file)
	stories, err := story.Unmarshal(json)

	if err != nil {
		fmt.Println("Unable to marshal stories JSON")
		os.Exit(1)
	}

	// in memory map of stories
	storyRepo := story.NewStoryRepo(stories)

	http.HandleFunc("/", serverStory(storyRepo, *templateName))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// Serve the story based on the URL path
func serverStory(repo *story.Repo, layout string) http.HandlerFunc {

	// closes over repo
	return func(writer http.ResponseWriter, request *http.Request) {
		p := strings.TrimPrefix(request.URL.Path, "/")

		if p == "" {
			http.Redirect(writer, request, "/intro", http.StatusMovedPermanently)
			return
		}

		s, err := repo.GetStory(p)

		if err != nil || s == nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		err = template.ServeStoryTemplate(writer, s, layout)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
