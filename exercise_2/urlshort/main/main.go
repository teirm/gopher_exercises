package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/teirm/gopher_exercises/exercise_2/urlshort"
)

// Read a file specified by the given path
func readFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func main() {

	jsonFile := flag.String("json-file", "", "name of json file")
	yamlFile := flag.String("yaml-file", "", "name of yaml file")
	flag.Parse()

	if *jsonFile != "" && *yamlFile != "" {
		log.Fatal("Too many input. Specify either jsonFile or yamlFile")
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var handler http.HandlerFunc
	if *jsonFile != "" {
		jsn, err := readFile(*jsonFile)
		if err != nil {
			log.Fatal("Unable to read json file '%s': %v\n", *jsonFile, err)
		}
		handler, err = urlshort.JSONHandler(jsn, mapHandler)
		if err != nil {
			log.Fatal("Unable to create JSONHandler: %v\n", err)
		}
	} else if *yamlFile != "" {
		yml, err := readFile(*yamlFile)
		if err != nil {
			log.Fatal("Unable to read yaml file '%s': %v\n", err)
		}
		handler, err = urlshort.YAMLHandler(yml, mapHandler)
		if err != nil {
			log.Fatal("Unable to create YAMLHandler: %v\n", err)
		}
	} else {
		handler = mapHandler
	}
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	fmt.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", handler)
	fmt.Println("Exiting server: %v\n", err)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
