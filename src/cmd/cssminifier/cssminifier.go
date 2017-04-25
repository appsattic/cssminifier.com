package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gomiddleware/mux"
)

var dir = "/var/lib/com-cssminifier"

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
	dir := os.Getenv("CSSMINIFIER_DIR")
	if dir == "" {
		log.Fatal("Specify a dir to save files to in the environment variable 'CSSMINIFIER_DIR'")
	}

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

	// pages
	m.Get("/", servePage(tmpl, "index", baseUrl, googleAnalytics))
	m.Get("/plugins", servePage(tmpl, "plugins", baseUrl, googleAnalytics))
	m.Get("/programs", servePage(tmpl, "programs", baseUrl, googleAnalytics))
	m.Get("/wget", servePage(tmpl, "wget", baseUrl, googleAnalytics))
	m.Get("/curl", servePage(tmpl, "curl", baseUrl, googleAnalytics))
	m.Get("/nodejs", servePage(tmpl, "nodejs", baseUrl, googleAnalytics))
	m.Get("/python", servePage(tmpl, "python", baseUrl, googleAnalytics))
	m.Get("/ruby", servePage(tmpl, "ruby", baseUrl, googleAnalytics))
	m.Get("/perl", servePage(tmpl, "perl", baseUrl, googleAnalytics))
	m.Get("/php", servePage(tmpl, "php", baseUrl, googleAnalytics))
	m.Get("/c-sharp", servePage(tmpl, "c-sharp", baseUrl, googleAnalytics))

	m.Get("/raw", redirect("/"))
	m.Post("/raw", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/raw: entry\n")

		// get the input
		input := r.FormValue("input")

		// minify this input
		file, err := minifyFile(input)
		if err != nil {
			log.Printf("Error minifying file")
			internalServerError(w, err)
			return
		}

		// stream to the response
		count, err := io.Copy(w, file)
		if err != nil {
			log.Printf("Error copying the minified file to the response")
			internalServerError(w, err)
			return
		}
		fmt.Printf("Read %d bytes from minified file.\n", count)
	})

	// finally, check all routing was added correctly
	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
