package runtime

import (
    "log"
    "strings"
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
            c.crunch(task.Id, true)
        }
        // Minutes computing
        c.crunch(task.Id, false)
    }
}

// Crunches the data for the given minute.
func (c *Cruncher) crunch(t time.Time, hour bool) {
    dao := db.NewTweetDAO(c.App.Mongo)
    crunchType := "minute"
    if hour {
        crunchType = "hour"
    }
    log.Printf("[info] [crunch] [%s] Will crunch the minute : %s\n", crunchType, t)

    var tweets []db.Tweet
    var err error
    if hour {
        tweets, err = dao.GetHourBucket(t)
    } else {
        tweets, err = dao.GetMinuteBucket(t)
    }

    if err != nil {
        log.Printf("[err] [crunch] While retrieving the bucket for time %s : %s\n", crunchType, t, err.Error())
        return
    }

    log.Printf("[info] [crunch] [%s] Retrieved %d tweets to crunch\n", crunchType, len(tweets))

    // Do the math.
    err = c.aggregateTweets(tweets)

    if err != nil {
        log.Printf("[err] [crunch] [%s] Error while computing data for tweets of : %s\n", crunchType, t)
    }
}

func (c *Cruncher) aggregateTweets(tweets []db.Tweet) error {
    // TODO This method could be speed up by 
    // ordering the tweets with their query 
    // and using regexp for the query.

    // In this map, I'll store the number of occurences
    // of each word located just after the query
    var data map[string]int
    data = make(map[string]int)

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
    log.Println(data)

    return nil
}
