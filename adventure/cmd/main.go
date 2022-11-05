package main

import (
	"adventure/internal/story"
	"adventure/internal/template"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("gophers.json")
	if err != nil {
		fmt.Println("Unable to read gophers.json")
		os.Exit(1)
	}
	defer file.Close()

	json, err := ioutil.ReadAll(file)
	stories, err := story.Unmarshal(json)

	if err != nil {
		fmt.Println("Unable to marshal gophers JSON")
		os.Exit(1)
	}

	// in memory map of stories
	storyRepo := story.NewStoryRepo(stories)

	http.HandleFunc("/", serverStory(storyRepo))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// Serve the story based on the URL path
func serverStory(repo *story.Repo) http.HandlerFunc {

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

		err = template.ServeStoryTemplate(writer, s)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
