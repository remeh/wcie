package runtime

import (
    "log"
    "strings"
    "time"

//  "github.com/ChimeraCoder/anaconda"
    "github.com/remeh/wcie/db"
)

const (
    TASK_COUNT_DONE_EACH_ITERATION = 50
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
    taskDAO := db.NewCrunchingTaskDAO(c.App.Mongo)
    tasks, err := taskDAO.GetNext(TASK_COUNT_DONE_EACH_ITERATION)
    if err != nil {
        log.Printf("[error] While retrieving some tasks to do : %s\n", err.Error())
        return
    }

    crunched := make([]db.CrunchingTask, 0)

    // Look whether its a minute or an hour to compute.
    for _, task := range tasks {
        done := false
        // Special case for hours.
        if task.Id.Minute() == 0 {
            // Crunch as minute and hour
            done = (c.crunch(task.Id, false) && c.crunch(task.Id, true))
        } else {
            // Minutes computing
            done = c.crunch(task.Id, false)
        }

        if done {
            crunched = append(crunched, task)
        }
    }

    err = taskDAO.RemoveAll(crunched)
    if err != nil {
        log.Printf("[err] [crunch] Error while removing crunched tasks : %s\n", err.Error())
    }
}

// Crunches the data for the given minute.
// Returns whether or not this data has been crunched
func (c *Cruncher) crunch(t time.Time, hour bool) bool {
    dao := db.NewTweetDAO(c.App.Mongo)
    crunchType := "minute"
    if hour {
        crunchType = "hour"
    }

    // Checks if it's time to compute this bucket
    if !c.isOver(t, hour) {
        return false
    }

    log.Printf("[info] [crunch] [%s] Will crunch the %s : %s\n", crunchType, crunchType, t)

    var tweets []db.Tweet
    var err error
    if hour {
        tweets, err = dao.GetHourBucket(t)
    } else {
        tweets, err = dao.GetMinuteBucket(t)
    }

    if err != nil {
        log.Printf("[err] [crunch] While retrieving the bucket for time %s : %s\n", crunchType, t, err.Error())
        return false
    }

    log.Printf("[info] [crunch] [%s] Retrieved %d tweets to crunch\n", crunchType, len(tweets))

    // Do the math.
    err = c.aggregateTweets(tweets)

    if err != nil {
        log.Printf("[err] [crunch] [%s] Error while computing data for tweets of : %s\n", crunchType, t)
        return false
    }

    return true
}

// To be sure that the hour / minute is finished
func (c *Cruncher) isOver(t time.Time, hour bool) bool {
    var plusOne time.Time
    if hour {
        plusOne = t.Add(time.Duration(1)*time.Hour)
        plusOne = time.Date(plusOne.Year(), plusOne.Month(), plusOne.Day(), plusOne.Hour(), 0, 0, 0, t.Location())
    } else {
        plusOne = t.Add(time.Duration(1)*time.Minute)
        plusOne = time.Date(plusOne.Year(), plusOne.Month(), plusOne.Day(), plusOne.Hour(), plusOne.Minute(), 0, 0, t.Location())
    }

    // Test that it's time to compute it
    if plusOne.Before(time.Now()) {
        return true
    }
    return false
}

func (c *Cruncher) aggregateTweets(tweets []db.Tweet) error {
    // TODO This method could be speed up by 
    // ordering the tweets with their query 
    // and using regexp for the query.

    // In this map, I'll store the number of occurences
    // of each word located just after the query
    data := make(map[string]int)

    // This algorithm could be speed up
    // Too many Trim call
    for _, tweet := range tweets {
        parts := strings.Split(strings.ToLower(tweet.Text), tweet.Query)
        if len(parts) == 1 {
            log.Printf("[warn] [crunch] Unable to explode the query with id '%d', text '%s'\n", tweet.TweetId, tweet.Text)
            continue
        }
        trimmed := strings.Trim(parts[1], " .,!?")
        nextSpace := strings.Index(trimmed, " ")
        word := ""
        if nextSpace == -1 {
            word = trimmed
        } else {
            word = strings.Trim(trimmed[0:nextSpace], ".!?,")
        }

        if data[word] != 0 {
            data[word] = data[word] + 1
        } else {
            data[word] = 1
        }
    }

    // TODO save this data.

    return nil
}
