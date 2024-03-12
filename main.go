package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("web/index.html"))
        tmpl.Execute(w, nil)
    })
    http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Clicked!")
    })

    fmt.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
