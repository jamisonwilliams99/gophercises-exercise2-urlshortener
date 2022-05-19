package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercises/urlshort/students/jamisonwilliams99/urlshort"
)

func main() {

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	} else {
		var data string
		switch os.Args[1] {
		case "-y":
			data = getYaml()
			handler, err := urlshort.YAMLHandler([]byte(data), mapHandler)
			if err != nil {
				panic(err)
			}
			fmt.Println("Starting the server on :8080")
			http.ListenAndServe(":8080", handler)
		case "-j":
			fmt.Println("Retrieving .json data...")
			data = getJSON()
			handler, err := urlshort.JSONHandler([]byte(data), mapHandler)
			if err != nil {
				panic(err)
			}
			fmt.Println("Starting the server on :8080")
			http.ListenAndServe(":8080", handler)
		default:
			log.Fatal("unknown argument")
		}
	}
}

func getJSON() string {
	var json string

	if len(os.Args) < 3 {
		json = `
[
	{"path": "/urlshort", "url": "https://github.com/gophercises/urlshort"},
	{"path": "/urlshort-final", "url": "https://github.com/gophercises/urlshort/tree/solution"}
]		
`
	} else {
		if os.Args[2] == "-f" {
			f, err := os.ReadFile("facebook.json")
			if err != nil {
				log.Fatal(err)
			}
			json = string(f)
		} else {
			log.Fatal("Unknown flag")
		}
	}
	return json
}

func getYaml() string {
	var yaml string

	if len(os.Args) < 3 {
		yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	} else {
		if os.Args[2] == "-f" {
			f, err := os.ReadFile("facebook.yaml")
			if err != nil {
				log.Fatal(err)
			}
			yaml = string(f)
		} else {
			log.Fatal("Unknown flag")
		}
	}
	return yaml
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
