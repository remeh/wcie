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
        if task.Minute() == 0 {
        }
    }
}

func (c *Cruncher) crunchMinute(t time.Time) {
    // TODO
}

func (c *Cruncher) crunchHour(t time.Time) {
    // TODO
}
