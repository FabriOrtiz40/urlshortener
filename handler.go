package main

import (
	"net/http"
	"log"
	"gopkg.in/yaml.v2"
	"encoding/json"
	bolt "go.etcd.io/bbolt"

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

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var paths []pathURL
	err := json.Unmarshal(jsonData, &paths)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(paths)
	return MapHandler(pathMap, fallback), nil
}

func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var targetURL string
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("paths"))
			if bucket == nil {
				return nil // No bucket, no redirect
			}
			val := bucket.Get([]byte(path))
			if val != nil {
				targetURL = string(val)
			}
			return nil
		})

		if err != nil {
			log.Printf("Error al acceder a la DB: %v", err)
			fallback.ServeHTTP(w, r)
			return
		}

		if targetURL != "" {
			http.Redirect(w, r, targetURL, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}


