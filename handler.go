package urlshort

import (
	"bytes"
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
		for path, url := range pathsToUrls {
			if r.URL.Path == path {
				println("redirecting to", url)
				http.Redirect(w, r, url, 301)
			}
		}
		fallback.ServeHTTP(w, r)
	}
}

// FileHandler will parse the provided YAML or JSON and then return
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
func FileHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	fileType := "yaml"
	if bytes.Contains(data, []byte("{")) {
		fileType = "json"
	}

	urls, parseErr := parseFile(fileType, data)

	// call MapHandler
	return MapHandler(mapFromUrls(urls), fallback), parseErr
}

func parseFile(fileType string, data []byte) ([]urlObj, error) {

	urls := []urlObj{}
	var err error
	switch fileType {
	case "yaml":
		err = yaml.Unmarshal(data, &urls)
	default:
		err = json.Unmarshal(data, &urls)
	}

	if err != nil {
		return nil, err
	}

	return urls, nil
}

func mapFromUrls(urls []urlObj) map[string]string {
	//array to map:
	urlMap := make(map[string]string)

	for _, o := range urls {
		urlMap[o.Path] = o.URL
	}

	return urlMap
}

type urlObj struct {
	//  the Unmarshal function requires public Uppercase properties
	URL      string `yaml:"url,omitempty"` // tags
	Path     string `yaml:"path,omitempty"`
	OtherVal string `yaml:"-"` // just for fun
}
