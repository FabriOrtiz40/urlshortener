package main

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(data []byte) ([]pathURL, error) {
	var paths []pathURL
	err := yaml.Unmarshal(data, &paths)
	return paths, err
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		// No encontrado → usar fallback
		fallback.ServeHTTP(w, r)

	}
}

func buildMap(paths []pathURL) map[string]string{
	m := make(map[string]string)
	for _, p := range paths {
		m[p.Path] = p.URL
	}
	return m
}

func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(paths)
	return MapHandler(pathMap, fallback), nil
}

