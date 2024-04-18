package main

import (
	"fmt"
	"log"

	"time"
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
    Name, Duration string
    Id, Day, Month, Year int
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
        duration integer,
        date datetime
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

func (m *model) addWorkout(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    d, err := time.Parse(time.DateOnly, r.Form.Get("date"))
    if err != nil {
        log.Fatal("Invalid Date")
    }
    
    wo := workout{
        Name: r.Form.Get("name"),
        Duration: r.Form.Get("duration"),
        Month: int(d.Month()),
        Day: d.Day(),
        Year: d.Year(),
    }

    db := m.DB
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := tx.Prepare("insert into workouts(name, duration, date) values(?, ?, ?);")
    if err != nil {
        log.Fatal(err)
    }

    defer stmt.Close()

    _, err = stmt.Exec(wo.Name, wo.Duration, d)
    if err != nil {
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }

    var html string
    for _, wo := range m.allWorkouts() {
        html += fmt.Sprintf(
            `<tr>
                <th scope="row">%s</th>
                <td>%s</td>
                <td>%d/%d/%d</td>
                <td><i id="delete" class="fa-solid fa-trash" hx-delete="/workouts/%d"hx-target="#workouts"></i></td>
            </tr>`,
            wo.Name,
            wo.Duration,
            wo.Month,
            wo.Day,
            wo.Year,
            wo.Id,
        )
    }
    fmt.Fprint(w, html)
}

func (m *model) allWorkouts() []workout {
    db := m.DB
    rows, err := db.Query("select id, name, duration, date from workouts order by date desc;")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var workouts []workout
    for rows.Next() {
        var w workout
        var d time.Time

        err = rows.Scan(&w.Id, &w.Name, &w.Duration, &d)
        if err != nil {
	    log.Fatal(err)
        }
        
        w.Month = int(d.Month())
        w.Day = d.Day()
        w.Year = d.Year()
        workouts = append(workouts, w)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    return workouts
}

func (m *model) deleteWorkout(w http.ResponseWriter, r *http.Request, id string) {

    db := m.DB
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := tx.Prepare("delete from workouts where id=?;")
    if err != nil {
        log.Fatal(err)
    }

    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }

    var html string
    for _, wo := range m.allWorkouts() {
        html += fmt.Sprintf(
            `<tr>
                <th scope="row">%s</th>
                <td>%s</td>
                <td>%d/%d/%d</td>
                <td><i id="delete" class="fa-solid fa-trash" hx-delete="/workouts/%d" hx-target="#workouts"></i></td>
            </tr>`,
            wo.Name,
            wo.Duration,
            wo.Month,
            wo.Day,
            wo.Year,
            wo.Id,
        )
    }
    fmt.Fprint(w, html)
}

func deleteHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        wid := r.URL.Path[len("/workouts/"):]
        fn(w, r, wid)
    }
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
        tmpl.Execute(w, m.allWorkouts())
    })

    http.HandleFunc("POST /workouts", m.addWorkout)
    http.HandleFunc("DELETE /workouts/", deleteHandler(m.deleteWorkout))

    fmt.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
