package db

import (
	"errors"

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

func (d *TweetDAO) Save(tweet *Tweet) error {
	if len(tweet.Id) > 0 {
		return d.collection.Update(bson.M{"_id": tweet.Id}, tweet)
	} else {
		return d.collection.Insert(tweet)
	}
}

// Gets many tweets using the given time bucket.
// A time bucket represents a set of tweets regrouped by minute.
func (d *TweetDAO) GetMinuteBucket(bucket time.Time) ([]Tweet, error) {
	// Retrieves some questions
	var q *mgo.Query
	// Do we allow already used questions ?
	if used {
		// TODO avoid to take too recent questions
		q = d.collection.Find(bson.M{"minute_bucket": bucket})
	}

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
