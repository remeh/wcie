package db

import (
    "time"

	"gopkg.in/v2/mgo"
	"gopkg.in/v2/mgo/bson"
)

// ----------------------
// Declarations

// DAO for tweet collection.
// @author Rémy MATHIEU
type TweetDAO struct {
	mongo      *Mongo
	collection *mgo.Collection
}

// ----------------------
// Methods

func NewTweetDAO(m *Mongo) *TweetDAO {
	return &TweetDAO{m, m.GetCollection(C_TWEET)}
}

// Save/update the given tweet.
func (d *TweetDAO) Upsert(tweet *Tweet) error {
	if len(tweet.Id) > 0 {
		return d.collection.Update(bson.M{"_id": tweet.Id}, tweet)
	} else {
		return d.collection.Insert(tweet)
	}
}

// Finds tweets by the tweet id.
func (d *TweetDAO) FindByTweetId(id int64) ([]Tweet, error) {
	var q *mgo.Query

	q = d.collection.Find(bson.M{"tweet_id": id})
    return d.unrollQuery(q)
}

// Gets many tweets using the given time bucket.
// A time bucket represents a set of tweets grouped by minute.
func (d *TweetDAO) GetHourBucket(bucket time.Time) ([]Tweet, error) {
	var q *mgo.Query

	q = d.collection.Find(bson.M{"hour_bucket": bucket})
    return d.unrollQuery(q)
}

// Gets many tweets using the given time bucket.
// A time bucket represents a set of tweets grouped by minute.
func (d *TweetDAO) GetMinuteBucket(bucket time.Time) ([]Tweet, error) {
	var q *mgo.Query

	q = d.collection.Find(bson.M{"minute_bucket": bucket})
    return d.unrollQuery(q)
}

// Read the whole iterator to return tweets.
func (d *TweetDAO) unrollQuery(q *mgo.Query) ([]Tweet, error) {
	iter := q.Iter()
	var tweet Tweet
	tweets := make([]Tweet, 0)
	for iter.Next(&tweet) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return tweets, nil
}
