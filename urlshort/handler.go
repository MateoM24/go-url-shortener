package urlshort

import (
	"fmt"
	"log"
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
		requestedURI := *&r.RequestURI
		if u, ok := pathsToUrls[requestedURI]; ok {
			http.Redirect(w, r, u, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
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
	parsedYaml := parseYaml(yml)
	pathMap := buildMap(parsedYaml)
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), nil
}

func parseYaml(yml []byte) []UrlPathPair {
	fmt.Println(yml)
	result := new([]UrlPathPair)
	err := yaml.Unmarshal(yml, &result)
	if err != nil {
		log.Fatalln("Cannot parse yaml config file")
	}
	fmt.Println(*result)
	return *result
}

func buildMap(pairs []UrlPathPair) map[string]string {
	configmap := make(map[string]string)
	for _, p := range pairs {
		configmap[p.Path] = p.URL
	}
	return configmap
}

type UrlPathPair struct {
	Path string
	URL  string
}
