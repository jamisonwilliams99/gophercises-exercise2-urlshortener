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
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
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
	// 1. Parse the yaml
	pathUrls, err := YamlParser(yml)
	if err != nil {
		return nil, err
	}

	// 2. Convert yaml array into map
	pathsToUrls := ConvYmlArrayToMap(pathUrls)

	// 3. return a map handler using the map
	return MapHandler(pathsToUrls, fallback), nil
}

func YamlParser(yml []byte) ([]yamlPathUrl, error) {
	var pathUrls []yamlPathUrl
	err := yaml.Unmarshal(yml, &pathUrls) // unmarshal will take yaml data and convert it to the struct that we defined (there are multiple yaml objects, so it is converted to a list of the structs)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func ConvYmlArrayToMap(pathUrls []yamlPathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

type yamlPathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 1. Parse the json
	pathUrls, err := JSONParser(jsn)
	if err != nil {
		return nil, err
	}

	// 2. Convert json array into map
	pathsToUrls := ConvJSONArrayToMap(pathUrls)

	// 3. return a map handler using the map
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONParser(jsn []byte) ([]jsonPathUrl, error) {
	var pathUrls []jsonPathUrl
	err := json.Unmarshal(jsn, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func ConvJSONArrayToMap(pathUrls []jsonPathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

type jsonPathUrl struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}
