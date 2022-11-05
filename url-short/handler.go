package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
	"urlshort/pkg/urlstore"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		path := normalisePath(request.URL.Path)

		url, hasKey := pathsToUrls[path]
		if !hasKey {
			fallback.ServeHTTP(writer, request)
			return
		}

		http.RedirectHandler(url, http.StatusPermanentRedirect).ServeHTTP(writer, request)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	links := make([]urlstore.Redirect, 0)
	err := yaml.Unmarshal(yml, &links)

	if err != nil {
		return nil, err
	}

	linkMap := createLinkMap(links)
	return MapHandler(linkMap, fallback), nil
}

// JsonHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//		[
//			{
//				"path": "/bbc",
//				"url": "https://www.bbc.co.uk/"
//			},
//			{
//				"path": "/mail",
//				"url": "https://mail.google.com/"
//			}
//		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JsonHandler(rawJson []byte, fallback http.Handler) (http.HandlerFunc, error) {
	links := make([]urlstore.Redirect, 0)
	err := json.Unmarshal(rawJson, &links)
	if err != nil {
		return nil, err
	}

	linkMap := createLinkMap(links)
	return MapHandler(linkMap, fallback), nil
}

func DbHandler(store urlstore.Store, fallback http.Handler) (http.HandlerFunc, error) {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := normalisePath(request.URL.Path)

		link, err := store.Get(path)

		if err != nil || link == nil {
			fallback.ServeHTTP(writer, request)
			return
		}

		http.RedirectHandler(link.Url, http.StatusPermanentRedirect).ServeHTTP(writer, request)

	}, nil
}

func createLinkMap(redirects []urlstore.Redirect) map[string]string {
	response := make(map[string]string)
	for _, redirect := range redirects {
		response[redirect.Path] = redirect.Url
	}

	return response
}

func normalisePath(path string) string {

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}
