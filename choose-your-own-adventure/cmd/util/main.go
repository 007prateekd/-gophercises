package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text    string `json:"text"`
		Chapter string `json:"arc"`
	} `json:"options"`
}

func ParseJson(fileBytes []byte) (Story, error) {
	story := make(Story)
	err := json.Unmarshal(fileBytes, &story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

var (
	DefaultHandlerTemp string
	DefaultTemplate    *template.Template
)

func InitializeTemplate(path string) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	DefaultHandlerTemp = string(fileBytes)
	DefaultTemplate = template.Must(template.New("").Parse(DefaultHandlerTemp))
}

type handler struct {
	story  Story
	tmpl   *template.Template
	pathFn func(r *http.Request) string
}

const (
	ERR_INTERNAL_SERVER   = "Something Went Wrong"
	ERR_CHAPTER_NOT_FOUND = "Chapter Not Found"
)

func DefaultPathFn(r *http.Request) string {
	path := r.URL.Path
	key := "intro"
	if path != "" && path != "/" {
		key = path[1:]
	}
	return key
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.pathFn(r)
	if chapter, ok := h.story[key]; ok {
		err := h.tmpl.Execute(w, chapter)
		if err != nil {
			log.Println(err)
			http.Error(w, ERR_INTERNAL_SERVER, http.StatusInternalServerError)
		}
	} else {
		http.Error(w, ERR_CHAPTER_NOT_FOUND, http.StatusNotFound)
	}
}

type HandlerOpt func(h *handler)

// Functional Option(s)
func Template(t *template.Template) HandlerOpt {
	return func(h *handler) {
		h.tmpl = t
	}
}

func PathFn(f func(r *http.Request) string) HandlerOpt {
	return func(h *handler) {
		h.pathFn = f
	}
}

func NewHandler(s Story, opts ...HandlerOpt) http.Handler {
	h := handler{story: s}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}
