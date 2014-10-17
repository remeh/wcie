package runtime

import (
    "log"
    "time"

    "github.com/remeh/wcie/crawl"
)

const (
    FREQUENCY = 10 // Job frequency of execution. (unit: minutes)
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

    log.Printf("Launching the crawling job, will execute every %d minutes.\n", FREQUENCY)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("Executing the crawling job.")

            // Generates for the next days
            crawl.Crawl()

            nextExecution = a.programNextCrawl()
        }

        time.Sleep(time.Minute * 1)
    }
}

// Programs the next execution of the job.
func (a *App) programNextCrawl() time.Time {
    now := time.Now()
    duration := time.Minute * FREQUENCY
    return now.Add(duration)
}
