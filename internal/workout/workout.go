package workout

import (
    "log"
    "time"
    "database/sql"
)

type Workout struct {
    Name string
    Duration int
}

func Add(db *sql.DB, w Workout) {
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    stmt, err := tx.Prepare("insert into workouts(name, date, duration) values(?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }

    defer stmt.Close()

    _, err = stmt.Exec(w.Name, time.Now().Local(), w.Duration)
    if err != nil {
        log.Fatal(err)
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}

