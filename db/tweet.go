package db

import (
    "log"
    "time"

    "gopkg.in/v2/mgo/bson"
    "github.com/ChimeraCoder/anaconda"
)

// ---------------------- 

// A crawled tweet
// @author Rémy MATHIEU
type Tweet struct {
    Id bson.ObjectId `bson:"_id,omitempty"`

    // The tweet id coming from Twitter
    TweetId int64 `bson:"tweet_id"`

    User User `bson:"user"`

    // The tweet content
    Text string `bson:"text"`

    // Which query has be done to retrieve this tweet
    Query string `bson:"query"`

    // Time at which this tweet has been crawled by the system
    CrawlingTime time.Time `bson:"crawling_time"`

    // Date at which the tweet has been published by the author.
    TweetTime time.Time `bson:"tweet_time"`

    // This bucket is used to regroup tweets by minute. It is based on the tweet time
    // and the seconds value is always set to 0 before insertion.
    MinuteBucket time.Time `bson:"minute_time_bucket"`

    // This bucket is used to regroup collected tweets by hour. It is based on the tweet time
    // and the minutes and seconds value is always set to 0 before insertion.
    HourBucket time.Time `bson:"hour_time_bucket"`
}

// Creates a new Tweet from an anaconda.Tweet
func NewTweetFromApiTweet(tweet *anaconda.Tweet, query string) *Tweet {
    // Some times handling
    now := time.Now()
    tweetTime, err := tweet.CreatedAtTime()
    if err != nil {
        log.Println("[warn] Unable to used tweet.CreatedAtTime() for tweet " + tweet.IdStr + ": using now as value.")
        tweetTime = now
    }

    minuteBucket := time.Date(tweetTime.Year(), tweetTime.Month(), tweetTime.Day(), tweetTime.Hour(), tweetTime.Minute(), 0, 0, tweetTime.Location())
    // Note that the hour bucket has 1 for second, it's a way to differentiate it from minute bucket
    hourBucket := time.Date(tweetTime.Year(), tweetTime.Month(), tweetTime.Day(), tweetTime.Hour(), 0, 1, 0, tweetTime.Location())

    // Creates the tweet.
    return &Tweet{
        TweetId: tweet.Id,
        User: User{
            Id: tweet.User.Id,
            ScreenName: tweet.User.ScreenName,
        },
        Text: tweet.Text,
        Query: query,
        CrawlingTime: now,
        TweetTime: tweetTime,
        MinuteBucket: minuteBucket,
        HourBucket: hourBucket,
    }
}
