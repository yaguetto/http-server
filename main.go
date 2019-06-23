package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	htmlDir = "/var/www/html"
)

func main() {
	handler()
	handlerEn()
	handlerPt()
	handleStatic()
	var port = ":8080"
	log.Printf("Listen and serve on %s", port)
	http.ListenAndServe(port, nil)
}

func handler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New requisition")
		checkLanguage(w, r)
	})
}

// can serve static files
func handleStatic() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/var/www/static"))))
}

// redirects based on the browser's language
func checkLanguage(w http.ResponseWriter, r *http.Request) {
	var lang = r.Header.Get("Accept-Language")
	if strings.HasPrefix(lang, "pt") {
		log.Printf("Requisition lang = %s, redirected to /pt", lang)
		http.Redirect(w, r, "/pt", 302)
	} else {
		log.Printf("Requisition lang = %s, redirected to /en", lang)
		http.Redirect(w, r, "/en", 302)
	}
}

func handlerPt() {
	http.HandleFunc("/pt", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New requisition on /pt")
		var file string
		file = fmt.Sprintf("%s/indice.html", htmlDir)
		var content, template = getHTML(file)
		template.Execute(w, content)
	})
}

func handlerEn() {
	http.HandleFunc("/en", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New requisition on /en")
		var file string
		file = fmt.Sprintf("%s/index.html", htmlDir)
		var content, template = getHTML(file)
		template.Execute(w, content)
	})
}

func getHTML(file string) ([]byte, *template.Template) {
	var fileContent, err = ioutil.ReadFile(file)
	if err != nil {
		log.Panicf(fmt.Sprintf("ioutil.ReadFile(): %s", err))
	}
	var template *template.Template
	template, err = template.ParseFiles(file)
	if err != nil {
		log.Println(fmt.Sprintf("template.ParseFiles(): %s", err))
	}
	return fileContent, template
}
