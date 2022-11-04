package main

import (
	"fmt"
	"net/http"
	"os"
	"urlshort"
	"urlshort/pkg/urlstore"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := `
[
	{
		"path": "/bbc",
		"url": "https://www.bbc.co.uk/"
	},
	{
		"path": "/mail",
		"url": "https://mail.google.com/"
	}
]
`
	jsonHandler, err := urlshort.JsonHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	store := urlstore.NewUrlStore()
	dbErr := store.Put(&urlstore.Redirect{
		Path: "/weather",
		Url:  "https://www.metoffice.gov.uk/",
	})

	if dbErr != nil {
		fmt.Printf("There was an error setting up the database %s\n", dbErr)
		os.Exit(1)
	}

	dbHandler, err := urlshort.DbHandler(store, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
