package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":7070", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user-id")

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "user-id",
			Value: "xxxx",
		}
	}

	http.SetCookie(w, cookie)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
