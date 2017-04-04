package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/chilts/sid"
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

	m.Get("/raw", redirect("/"))
	m.Post("/raw", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/raw: entry\n")

		// get the input
		input := r.FormValue("input")

		// create a unique id to use in the filenames
		id := sid.Id()
		filename := path.Join(dir, id+".css")
		output := path.Join(dir, id+".min.css")

		fmt.Printf("id=%s\n", id)
		fmt.Printf("filename=%s\n", filename)
		fmt.Printf("output=%s\n", output)

		// write to a file
		fOrig, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			internalServerError(w, err)
			return
		}
		defer fOrig.Close()

		// copy from the input to the file
		n, err := io.WriteString(fOrig, input)
		if err != nil {
			internalServerError(w, err)
			return
		}
		fmt.Printf("Written %d bytes to original file.\n", n)

		// run `cleancss`
		cmd := exec.Command("./node_modules/.bin/cleancss", "--output", output, filename)
		err = cmd.Run()
		if err != nil {
			internalServerError(w, err)
			return
		}

		// open the file again and stream to the response
		fMin, err := os.Open(output)
		if err != nil {
			internalServerError(w, err)
			return
		}
		defer fMin.Close()

		// stream to the response
		n2, err := io.Copy(w, fMin)
		if err != nil {
			internalServerError(w, err)
			return
		}
		fmt.Printf("Read %d bytes from minified file.\n", n2)
	})

	// finally, check all routing was added correctly
	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
