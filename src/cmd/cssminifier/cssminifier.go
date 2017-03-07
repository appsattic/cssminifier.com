package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomiddleware/mux"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// setup
	baseUrl := os.Getenv("CSSMINIFIER_BASE_URL")
	port := os.Getenv("CSSMINIFIER_PORT")
	if port == "" {
		log.Fatal("Specify a port to listen on in the environment variable 'CSSMINIFIER_PORT'")
	}
	googleAnalytics := os.Getenv("CSSMINIFIER_GOOGLE_ANALYTICS")

	// load up all templates
	tmpl, err := template.New("").ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// the mux
	m := mux.New()

	// do some static routes before doing logging
	m.All("/s", fileServer("static"))
	m.Get("/favicon.ico", serveFile("./static/favicon.ico"))
	m.Get("/robots.txt", serveFile("./static/robots.txt"))

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			BaseUrl         string
			GoogleAnalytics string
		}{
			baseUrl,
			googleAnalytics,
		}
		render(w, tmpl, "index.html", data)
	})

	// finally, check all routing was added correctly
	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
