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
        // Special case for hours, the nano is to 1.
        if task.Id.Minute() == 0 && task.Id.Second() == 1 {
            // Crunch as minute and hour
            done = c.crunch(task.Id, true)
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
    tweetDAO := db.NewTweetDAO(c.App.Mongo)
    resultDAO := db.NewCruncherResultDAO(c.App.Mongo)

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
        tweets, err = tweetDAO.GetHourBucket(t)
    } else {
        tweets, err = tweetDAO.GetMinuteBucket(t)
    }

    if err != nil {
        log.Printf("[err] [crunch] While retrieving the bucket for time %s : %s\n", crunchType, t, err.Error())
        return false
    }

    log.Printf("[info] [crunch] [%s] Retrieved %d tweets to crunch\n", crunchType, len(tweets))

    // Do the math.
    err = c.aggregateTweets(resultDAO, t, tweets)

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

func (c *Cruncher) aggregateTweets(resultDAO *db.CruncherResultDAO, bucket time.Time, tweets []db.Tweet) error {
    // TODO This method could be speed up by 
    // ordering the tweets with their query 
    // and using regexp for the query.

    // In this map, I'll store the number of occurences
    // of each word located just after the query
    data := make(map[string]int)

    for _, tweet := range tweets {
        parts := strings.Split(strings.ToLower(tweet.Text), tweet.Query)
        if len(parts) == 1 {
            log.Printf("[warn] [crunch] Unable to explode the query with id '%d', text '%s'\n", tweet.TweetId, tweet.Text)
            continue
        }
        trimmed := strings.Trim(parts[1], " .,!?")
        cleaned := c.cleanWord(trimmed)
        nextSpace := strings.Index(cleaned, " ")
        word := ""
        if nextSpace == -1 {
            word = cleaned
        } else {
            word = cleaned[0:nextSpace]
        }

        if data[word] != 0 {
            data[word] = data[word] + 1
        } else {
            data[word] = 1
        }
    }

    return resultDAO.Upsert(&db.CruncherResult{Id: bucket, Data: data})
}

// Temporary and ugly method to clean the word.
func (c *Cruncher) cleanWord(word string) string {
    w1 := strings.Replace(word, ".", " ", -1)
    w2 := strings.Replace(w1, "!", " ", -1)
    w3 := strings.Replace(w2, "?", " ", -1)
    w4 := strings.Replace(w3, ".", " ", -1)
    w5 := strings.Replace(w4, "#", " ", -1)
    return w5
}
