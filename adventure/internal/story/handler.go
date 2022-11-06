package story

import (
	"net/http"
	"strings"
)

type handler struct {
	repo   *Repo
	layout string
	parser func(r *http.Request) string
}

// HandlerOption https://golang.cafe/blog/golang-functional-options-pattern.html
type HandlerOption func(*handler)

func NewHandler(repo *Repo, layout string, options ...HandlerOption) http.Handler {
	handler := handler{
		repo:   repo,
		layout: layout,
		parser: defaultRouter,
	}

	for _, opt := range options {
		opt(&handler)
	}

	return handler
}

// WithRouterOption Define a customer router function
// the function signature is func(r *http.Request) string.
// This function will take the request, and return a string representing the name of the story.
// The story name for example could be sourced from the request path, a header, or the body.
//
// If no router function is provided then the default router will be used.
// The default router assumes that the request URL path is the name of the story e.g. /hello where hello is the name of the story.
func WithRouterOption(parser func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.parser = parser
	}
}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	p := h.parser(request)

	s, err := h.repo.GetStory(p)

	if err != nil || s == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	err = serveStoryTemplate(writer, s, h.layout)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func defaultRouter(request *http.Request) string {
	p := strings.TrimPrefix(request.URL.Path, "/")

	if p == "" {
		return "intro"
	}

	return p
}
