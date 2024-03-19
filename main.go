package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"html/template"
    "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./ptlog.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    stmt := `
    create table if not exists workouts(
        id interger not null primary key,
        date datetime, name text,
        duration integer
    );
    delete from workouts;
    `
    _, err = db.Exec(stmt)
    if err != nil {
        log.Printf("%q: %s\n", err, stmt)
        return
    }

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
