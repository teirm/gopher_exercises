package urlshort

import (
	"encoding/json"
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

// Parsing entry for yaml files
type yamlEntry struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// Parse a byte array representing a yaml file
// return an array of structs or an error
func parseYaml(yml []byte) ([]yamlEntry, error) {
	entries := []yamlEntry{}
	if err := yaml.Unmarshal(yml, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Convert an array of Entry type structs into
// a mapping
func buildMapYaml(entries []yamlEntry) map[string]string {
	pathMap := make(map[string]string)
	for _, entry := range entries {
		pathMap[entry.Path] = entry.URL
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
	pathMap := buildMapYaml(entries)
	return MapHandler(pathMap, fallback), nil
}

// entry for json file parsing
type jsonEntry struct {
	Path string
	URL  string
}

// Parse a JSON file into an array of jsonEntry
func parseJSON(jsn []byte) ([]jsonEntry, error) {
	var entries []jsonEntry
	if err := json.Unmarshal(jsn, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// build a map from jsonEntry array
func buildMapJson(entries []jsonEntry) map[string]string {
	pathMap := make(map[string]string)
	for _, entry := range entries {
		pathMap[entry.Path] = entry.URL
	}
	return pathMap
}

// JSONHandler will parse the provided JSON and then
// return an http.HandlerFunc mapping paths to corresponding
// URL
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	pathMap := buildMapJson(entries)
	return MapHandler(pathMap, fallback), nil
}
