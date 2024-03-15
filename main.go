package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"html/template"
)

func main() {
    // Handle static files
    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("web/index.html"))
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("GET /clicked", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Clicked!")
    })

    fmt.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
