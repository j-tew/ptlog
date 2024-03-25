package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "time"
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type model struct {
    DB *sql.DB
    workout workout
}

type workout struct {
    Name string
    Day, Month, Year int
    Duration int
}

func (m *model) setup() error {
    db, err := sql.Open("sqlite3", "./ptlog.db")
    if err != nil {
        return err
    }

    stmt := `
    create table if not exists workouts(
        id integer primary key autoincrement,
        name text,
        day integer,
        month integer,
        year integer,
        duration integer
    );
    delete from workouts;
    `
    _, err = db.Exec(stmt)
    if err != nil {
        log.Printf("%q: %s\n", err, stmt)
        return err
    }
    m.DB = db
    return nil
}

func (m *model) addWorkout(w http.ResponseWriter, r *http.Request)  {
    var wo workout

    err := json.NewDecoder(r.Body).Decode(&wo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        log.Println(r.Body)
        return
    }

    db := m.DB
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := tx.Prepare("insert into workouts(name, day, month, year, duration) values(?, ?, ?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }

    defer stmt.Close()

    _, err = stmt.Exec(wo.Name, wo.Day, wo.Month, wo.Year, wo.Duration)
    if err != nil {
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}

func (m *model) getWorkouts() []workout {
    db := m.DB
    rows, err := db.Query("select name, day, month, year, duration from workouts")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var workouts []workout
    for rows.Next() {
        var w workout
        err = rows.Scan(&w.Name, &w.Day, &w.Month, &w.Year, &w.Duration)
        if err != nil {
	    log.Fatal(err)
        }
        workouts = append(workouts, w)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    return workouts
}

func main() {
    m := model{}

    err := m.setup()
    if err != nil {
        log.Fatal(err)
    }
    defer m.DB.Close()

    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
        tmpl := template.Must(template.ParseFiles("web/index.html"))
        tmpl.Execute(w, m.getWorkouts())
    })

    http.HandleFunc("POST /workout", m.addWorkout)

    fmt.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
