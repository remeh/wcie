package db

import (
    "time"

    "gopkg.in/v2/mgo/bson"
)

// ---------------------- 

// A crawled tweet
// @author Rémy MATHIEU
type Tweet struct {
    Id bson.ObjectId `bson:"_id,omitempty"`

    // The tweet id coming from Twitter
    TweetId string `bson:"tweet_id"`

    // The author id of the tweet. Example : 345678
    // XXX Is there really this data ?
    AuthorId int `bson:"author_id"`

    // The author of the tweet. Example : remeh
    Author string `bson:"author"`

    // The tweet
    Content string `bson:"content"`

    // Time at which this tweet has been crawled by the system
    CrawlingTime time.Time `bson:"crawling_time"`

    // Date at which the tweet has been published by the author.
    TweetTime time.Time `bson:"tweet_time"`

    // This bucket is used to regroup tweets by minute. It iss base on the tweet time
    // and the seconds value is always set to 0 before insertion.
    MinuteBucket time.Time `bson:"minute_time"`

    // This bucket is used to regroup tweets by hour. It iss base on the tweet time
    // and the minutes and seconds value is always set to 0 before insertion.
    HourBucket time.Time `bson:"hour_time"`

}

