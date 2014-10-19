package runtime

import (
    "log"
    "time"

//  "github.com/ChimeraCoder/anaconda"
    "github.com/remeh/wcie/db"
)

// Our crawler.
type Cruncher struct {
    App *App
}

func NewCruncher(app *App) *Cruncher {
    return &Cruncher{App: app}
}

// Takes some task to do and crunch the data
func (c *Cruncher) Crunch() {
    tasks, err := db.NewCrunchingTaskDAO(c.App.Mongo).GetNext(5)
    if err != nil {
        log.Printf("[error] While retrieving some tasks to do : %s\n", err.Error())
        return
    }

    // Look whether its a minute or an hour to compute.
    for _, task := range tasks {
        // Special case for hours.
        if task.Id.Minute() == 0 {
            c.crunchHour(task.Id)
        }
        // Minutes computing
        c.crunchMinute(task.Id)
    }
}

// Crunches the data for the given minute.
func (c *Cruncher) crunchMinute(t time.Time) {
    dao := db.NewTweetDAO(c.App.Mongo)
    log.Printf("[info] [crunch] Will crunch the minute : %s\n", t)

    tweets, err := dao.GetMinuteBucket(t)

    if err != nil {
        log.Printf("[err] [crunch] While retrieving the bucket for time %s : %s\n", t, err.Error())
        return
    }

    log.Printf("[info] [crunch] Retrieved %d tweets to crunch\n", len(tweets))
}

// Crunches the data for the given hour.
func (c *Cruncher) crunchHour(t time.Time) {
    // TODO
}
