package main

import (
	"fmt"
	"io"
	"log"
    "time"
	"net/http"
	"html/template"
    "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func addWorkout(db *sql.DB, name string, duration int) {
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := tx.Prepare("insert into workouts(name, date, duration) values(?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }

    defer stmt.Close()

    _, err = stmt.Exec(name, time.Now().Local(), duration)
    if err != nil {
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    db, err := sql.Open("sqlite3", "./ptlog.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    stmt := `
    create table if not exists workouts(
        id integer primary key autoincrement,
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

    addWorkout(db, "run", 30)

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
