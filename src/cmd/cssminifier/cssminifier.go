package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

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

		// run `cleancss` (shouldn't be anything on stdout, but perhaps stderr if there's a problem)
		cmd := exec.Command("./node_modules/.bin/cleancss", "--output", output, filename)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			internalServerError(w, err)
			return
		}

		// start the process
		err = cmd.Start()
		if err != nil {
			internalServerError(w, err)
			return
		}

		// slurp in stderr
		b, _ := ioutil.ReadAll(stderr)
		slurpStdErr := string(b)
		fmt.Printf("stderr : %s\n", slurpStdErr)

		// wait for the command to finish
		err = cmd.Wait()
		if err != nil {
			internalServerError(w, err)
			return
		}

		if len(slurpStdErr) > 0 {
			fmt.Printf("Something appeared on stderr = %d\n", len(slurpStdErr))
			// replace some bits we don't want to show the user (ToDo: generate this from `dir`.)
			r1 := strings.NewReplacer("../../../../var/lib/cssminifier/css/", "", "*/", "", "/*", "")
			slurpStdErr = r1.Replace(slurpStdErr)

			// wrap the error in a CSS comment at the top of the minified CSS so the user will notice it
			slurpStdErr = "/*\n\n" + slurpStdErr + "\n*/\n\n"
			fmt.Printf("stderr : %s\n", slurpStdErr)
		}

		// open the file again and stream to the response
		fMin, err := os.Open(output)
		if err != nil {
			internalServerError(w, err)
			return
		}
		defer fMin.Close()

		// write any error to the output first
		_, err = w.Write([]byte(slurpStdErr))
		if err != nil {
			internalServerError(w, err)
			return
		}

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
