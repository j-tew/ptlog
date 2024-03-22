# ptlog

A minimal web app to record physical activity.

Finally ready to mess with HTTP in Go and htmx. Starting out with a simple workout tracker.

## Notes

*To pass to template, struct fields must be exportable*

### So far:

- Basic template with one htmx get action that swaps the innerHTML of a button
- SQLite3 DB and methods to add a workout and fetch all workouts
- Some templating to pass the workouts to the UI (just the name for now)

### On the horizon:

- [ ] More methods to handle workouts (edit/remove)
- [ ] Render workouts in a table
- [ ] Allow adding workouts with form input (sanitize)
- [ ] Handle date and time for workouts more gracefully
