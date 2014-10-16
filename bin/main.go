// The Job crawling tweets.
package main

import (
    "log"
    "time"

    "github.com/remeh/wcie"
    "github.com/remeh/wcie/db"
)

const (
    FREQUENCY = 10 // Job frequency of execution. (unit: minutes)
)

func main() {
    nextExecution := time.Now()

    log.Printf("Launching the crawling job, will execute every %d minutes.\n", FREQUENCY)

    // Mainloop, yo.
    for {
        // Do we have to run ?
        if nextExecution.Before(time.Now()) {
            log.Println("Executing the crawling job.")

            // Generates for the next days
            wcie.Crawl()

            nextExecution = programNextExecution()
        }

        time.Sleep(time.Minute * 1)
    }
}

func generateNextDays() {
    // Gets a Mongo connection.
    m := db.GetConnection()
    defer m.Close()

    // TODO
}

// Programs the next execution of the job.
func programNextExecution() time.Time {
    now := time.Now()
    duration := time.Minute * FREQUENCY
    return now.Add(duration)
}
