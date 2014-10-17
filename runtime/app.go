package runtime

import (
    "log"
    "time"
)

const (
    FREQUENCY = 30 // Job frequency of execution. (unit: seconds). Can't be < 5
)

type App struct {
    // App configuration.
    Config Config
}


func NewApp(config Config) *App {
    return &App{Config: config}
}

// Starts the application
func (a *App) Start() {
    a.startCrawlJob()
}

// Starts the crawling Job.
func (a *App) startCrawlJob() {
    nextExecution := time.Now()

    log.Printf("Launching the crawling job, will execute every %d seconds.\n", FREQUENCY)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("Executing the crawling job.")

            // Generates for the next days
            Crawl(a)

            nextExecution = a.programNextCrawl()
        }

        time.Sleep(time.Second * 5)
    }
}

// Programs the next execution of the job.
func (a *App) programNextCrawl() time.Time {
    now := time.Now()
    duration := time.Second * FREQUENCY
    return now.Add(duration)
}
