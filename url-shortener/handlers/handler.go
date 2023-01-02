package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Item struct {
	Path string
	URL  string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func parseYaml(yamlData []byte) ([]Item, error) {
	var parsedYaml []Item
	err := yaml.Unmarshal(yamlData, &parsedYaml)
	if err != nil {
		return nil, err
	}
	return parsedYaml, nil
}

func buildMap(parsedYaml []Item) map[string]string {
	pathMap := make(map[string]string)
	for _, item := range parsedYaml {
		pathMap[item.Path] = item.URL
	}
	return pathMap
}

func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yamlData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseJson(jsonData []byte) ([]Item, error) {
	var parsedJson []Item
	err := json.Unmarshal(jsonData, &parsedJson)
	if err != nil {
		return nil, err
	}
	return parsedJson, nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(jsonData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}
