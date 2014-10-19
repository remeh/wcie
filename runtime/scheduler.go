package runtime

import (
    "log"
    "time"

    "github.com/remeh/wcie/db"
)

// Our crawler.
type Scheduler struct {
    App *App
}

func NewScheduler(app *App) *Scheduler {
    return &Scheduler{App: app}
}

// The job of the scheduler is to say to the
// cruncher which period of time he has to crunch data for.
// To do so, it has a table in the DB storing its state.
func (s *Scheduler) Schedule() {
    state := db.NewSchedulerStateDAO(s.App.Mongo).Get()

    s.scheduleMinute(state)
    // TODO scheduleHour
}

// Takes the last computed minute and generate tasks to reach
// the new one.
// If nothing has been computed (the LastComputedMinute is at 0), 
// we just compute the current hour to the current minute.
func (s *Scheduler) scheduleMinute(state *db.SchedulerState) {
    now := time.Now()
    startMinute := state.LastScheduledMinute
    if state.LastScheduledMinute.IsZero() {
        // Schedule the compute of each minute
        startMinute = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
    }

    // Creates all the task to compute
    lastScheduledMinute := s.createMinuteTaskUntilNow(&now, &startMinute)

    // Tasks created, we must update the scheduling state
    state.LastScheduledMinute = *lastScheduledMinute
    db.NewSchedulerStateDAO(s.App.Mongo).Save(state)
}

// Create tasks reaching now minute per minute using the given time.
// Returns the end of the interval of already scheduled minutes
func (s *Scheduler) createMinuteTaskUntilNow(now *time.Time, t *time.Time) *time.Time {
    for t.Before(*now) {
        // Creates the crunching task.
        err := s.createCrunchingTaskFor(t)
        if err != nil {
            log.Printf("[err] Unable to create the crunching task for time : %s\n", t)
            log.Printf("[err] Reason : %s\n", err.Error())
        }

        // Adds one minute
        *t = t.Add( time.Duration(1) * time.Minute )
    }
    return t
}

// Creates the crunching task for the given time.
func (s *Scheduler) createCrunchingTaskFor(t *time.Time) error {
    // Creating the ID representing the date
    var value int64
    value = int64(t.Year() * 100 * 100 * 100 * 100)
    value += int64(int(t.Month()) * 100 * 100 * 100)
    value += int64(t.Day() * 100 * 100)
    value += int64(t.Hour() * 100)
    value += int64(t.Minute())

    // Creates and saves the task
    task := &db.CrunchingTask{Id: value, CreationTime: time.Now()}
    log.Printf("[info] Created crunching task for : %d\n", value)
    return db.NewCrunchingTaskDAO(s.App.Mongo).Upsert(task)
}
