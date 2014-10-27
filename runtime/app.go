package runtime

import (
    "log"
    "time"
    "os"
    "os/signal"

    "github.com/remeh/wcie/db"
)

const (
    FREQUENCY_CRAWL = 60 // Job frequency of execution. (unit: seconds). Can't be < 5
    FREQUENCY_CRUNCH = 60 // Job frequency of execution. (unit: seconds). Can't be < 5
    FREQUENCY_SCHEDULING = 60 // Job frequency of execution. (unit: seconds). Can't be < 5
)

type App struct {
    // App configuration.
    Config Config
    Mongo *db.Mongo
}


func NewApp(config Config) *App {
    return &App{Config: config, Mongo: db.GetConnection(config.MongoURI)}
}

// Starts the application
func (a *App) Start() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    go a.startCrawlJob()
    go a.startCrunchJob()
    go a.startSchedulingJob()

    <-c
    log.Println("[info] Closing.")
}

// TODO Refactor these two jobs.

// Starts the scheduling Job.
func (a *App) startSchedulingJob() {
    nextExecution := time.Now()

    log.Printf("[info] Creating the scheduling job, will execute every %d seconds.\n", FREQUENCY_CRUNCH)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("[info] Executing the scheduling job.")

            scheduler := NewScheduler(a)
            scheduler.Schedule()

            nextExecution = a.programNextScheduling()
            log.Println("[info] End of the scheduling job.")
        }
        time.Sleep(time.Second * 5)
    }
}

// Programs the next execution of the job.
func (a *App) programNextScheduling() time.Time {
    now := time.Now()
    duration := time.Second * FREQUENCY_SCHEDULING
    return now.Add(duration)
}
// Starts the crunching Job.
func (a *App) startCrunchJob() {
    nextExecution := time.Now()

    log.Printf("[info] Creating the crunching job, will execute every %d seconds.\n", FREQUENCY_CRUNCH)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("[info] Executing the crunching job.")

            cruncher := NewCruncher(a)
            cruncher.Crunch()

            nextExecution = a.programNextCrawl()
            log.Println("[info] End of the crunching job.")
        }
        time.Sleep(time.Second * 5)
    }
}

// Programs the next execution of the job.
func (a *App) programNextCrunch() time.Time {
    now := time.Now()
    duration := time.Second * FREQUENCY_CRUNCH
    return now.Add(duration)
}

// Starts the crawling Job.
func (a *App) startCrawlJob() {
    nextExecution := time.Now()

    log.Printf("[info] Creating the crawling job, will execute every %d seconds.\n", FREQUENCY_CRAWL)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("[info] Executing the crawling job.")

            // Let's crawl.
            crawler := NewCrawler(a)
            crawler.Crawl()

            nextExecution = a.programNextCrawl()
            log.Println("[info] End of the crawling job.")
        }
        time.Sleep(time.Second * 5)
    }
}

// Programs the next execution of the job.
func (a *App) programNextCrawl() time.Time {
    now := time.Now()
    duration := time.Second * FREQUENCY_CRAWL
    return now.Add(duration)
}
