package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		request := r.URL.Path
		if redirect, ok := pathsToUrls[request]; !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, redirect, 307)
		}
	}
}

type Entry struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// Parse a byte array representing a yaml file
// return an array of structs or an error
func parseYaml(yml []byte) ([]Entry, error) {
	entries := []Entry{}
	if err := yaml.Unmarshal(yml, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Convert an array of Entry type structs into
// a mapping
func buildMap(entries []Entry) map[string]string {
	pathMap := make(map[string]string)

	for _, entry := range entries {
		pathMap[entry.Path] = entry.Url
	}
	return pathMap
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
	entries, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(entries)
	return MapHandler(pathMap, fallback), nil
}
