package main

import (
	"fmt"
	"net/http"

	"github.com/gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	//   otherval: ignoreme
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `

	json := `
	[{ 
		"path": "/urlshort",
		"url": "https://github.com/gophercises/urlshort"
	},
	{
  "path": "/urlshort-final",
		"url": "https://github.com/gophercises/urlshort/tree/solution"
	}
	]
  `

	fileHandler, err := urlshort.FileHandler([]byte(json), mapHandler)
	if err != nil {
		println("\nUH OH - FileHandler returned err\n!")
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", fileHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
