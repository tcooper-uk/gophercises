package story

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Rels struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Story struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Rels   `json:"options"`
}

type Repo struct {
	arcs map[string]Story
}

func LoadAllStoriesFromJson(f *os.File) (map[string]Story, error) {
	json, err := ioutil.ReadAll(f)
	stories, err := unmarshal(json)

	if err != nil {
		return stories, err
	}

	return stories, nil
}

func NewStoryRepo(data map[string]Story) *Repo {
	return &Repo{
		data,
	}
}

func (repo *Repo) GetStory(arc string) (*Story, error) {
	// find a story.
	s, found := repo.arcs[arc]

	if !found {
		return nil, errors.New("unable to find story")
	}

	return &s, nil
}

func unmarshal(data []byte) (map[string]Story, error) {
	var stories map[string]Story
	err := json.Unmarshal(data, &stories)

	if err != nil {
		return nil, err
	}

	return stories, nil
}
