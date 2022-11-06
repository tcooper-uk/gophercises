package main

import (
	"adventure/internal/story"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {

	inputFile := flag.String("input", "gophers.json", "The JSON input file containing the stories. Default is gophers.json")
	templateName := flag.String("template", "layout.html", "The name of the template to render. Default is layout.html")
	flag.Parse()

	file := loadFile(*inputFile)
	defer file.Close()

	// in memory map of stories
	storyRepo := loadStoryRepo(file)

	// start http server
	handler := story.NewHandler(storyRepo, *templateName)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}
}

func loadStoryRepo(f *os.File) *story.Repo {
	stories, err := story.LoadAllStoriesFromJson(f)

	if err != nil {
		fmt.Println("Unable to marshal stories JSON")
		os.Exit(1)
	}

	return story.NewStoryRepo(stories)
}

func loadFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Unable to read %s\n", file)
		os.Exit(1)
	}

	return file
}
