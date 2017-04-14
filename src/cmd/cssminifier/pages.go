package main

import (
	"html/template"
	"net/http"
)

func servePage(tmpl *template.Template, name, baseUrl, googleAnalytics string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			BaseUrl         string
			PageName        string
			GoogleAnalytics string
		}{
			baseUrl,
			name,
			googleAnalytics,
		}
		render(w, tmpl, name+".html", data)
	}
}
