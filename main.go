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

type model struct {
    DB *sql.DB
    workout workout
}

type workout struct {
    Name string
    day, month, year int
    duration int
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

func (m *model) addWorkout(w workout) error {
    db := m.DB
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare("insert into workouts(name, day, month, year, duration) values(?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }

    defer stmt.Close()

    _, err = stmt.Exec(w.Name, w.day, w.month, w.year, w.duration)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }
    return nil
}

func (m *model) getWorkouts() ([]workout, error) {
    db := m.DB
    rows, err := db.Query("select name, day, month, year, duration from workouts")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var workouts []workout
    for rows.Next() {
        var w workout
        err = rows.Scan(&w.Name, &w.day, &w.month, &w.year, &w.duration)
        if err != nil {
            return nil, err
        }
        workouts = append(workouts, w)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return workouts, nil
}

func main() {
    m := model{}

    err := m.setup()
    if err != nil {
        log.Fatal(err)
    }
    defer m.DB.Close()

    year, month, day := time.Now().Date()
    wo := workout{Name: "run", day: day, month: int(month), year: year, duration: 30}

    err = m.addWorkout(wo)
    if err != nil {
        log.Fatal(err)
    }

    workouts, err := m.getWorkouts()
    if err != nil {
        log.Fatal(err)
    }

    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
        tmpl := template.Must(template.ParseFiles("web/index.html"))
        tmpl.Execute(w, workouts)
    })

    http.HandleFunc("GET /clicked", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Clicked!")
    })

    fmt.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
