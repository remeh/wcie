package runtime

import (
    "log"
    "time"
    "os"
    "os/signal"

    "github.com/remeh/wcie/db"
)

const (
    FREQUENCY = 60 // Job frequency of execution. (unit: seconds). Can't be < 5
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

    <-c
    log.Println("[info] Closing.")
}

// Starts the crawling Job.
func (a *App) startCrawlJob() {
    nextExecution := time.Now()

    log.Printf("[info] Creating the crawling job, will execute every %d seconds.\n", FREQUENCY)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("[info] Executing the crawling job.")

            // Generates for the next days
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
    duration := time.Second * FREQUENCY
    return now.Add(duration)
}
