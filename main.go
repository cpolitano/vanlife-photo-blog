package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	c := getCookie(w, req)

	if req.Method == http.MethodPost {
		multipartFile, fileHeader, err := req.FormFile("new-file")
		if err != nil {
			fmt.Println(err)
		}
		defer multipartFile.Close()

		fileExtension := strings.Split(fileHeader.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, multipartFile)
		fileName := fmt.Sprintf("%x", h.Sum(nil)) + "." + fileExtension // print h as hexidecimal
		// create new file
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(workingDir, "public", "images", fileName)
		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer newFile.Close()

		multipartFile.Seek(0, 0)        // reset to beginning of file
		io.Copy(newFile, multipartFile) // copy file to public/images
		c = appendValue(w, c, fileName)
	}

	xs := strings.Split(c.Value, "|") // creates slice of strings
	tpl.ExecuteTemplate(w, "index.gohtml", xs)
}

func getCookie(w http.ResponseWriter, req *http.Request) *http.Cookie {
	cookie, err := req.Cookie("session")

	if err == http.ErrNoCookie {
		sessionId := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: sessionId.String(),
		}
		http.SetCookie(w, cookie)
	}
	return cookie
}

func appendValue(w http.ResponseWriter, c *http.Cookie, fname string) *http.Cookie {

	s := c.Value

	if !strings.Contains(s, fname) {
		s += "|" + fname
	}

	c.Value = s
	http.SetCookie(w, c)
	return c
}
