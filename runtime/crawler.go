package runtime

import (
    "log"
    "net/url"

    "github.com/ChimeraCoder/anaconda"
    "github.com/remeh/wcie/db"
)

// Our crawler.
type Crawler struct {
    App *App
}

func NewCrawler(app *App) *Crawler {
    return &Crawler{App: app}
}

// TODO cache tweet by ids.
func (c *Crawler) Crawl() {
    // Api twitter provided by ChimeraCoder !
    anaconda.SetConsumerKey(c.App.Config.TwitterApiKey)
    anaconda.SetConsumerSecret(c.App.Config.TwitterSecret)
    api := anaconda.NewTwitterApi(c.App.Config.TwitterAccessToken, c.App.Config.TwitterAccessTokenSecret)
    defer api.Close()

    c.Search(api, "\"je mange un\"")
    c.Search(api, "\"je mange une\"")
    c.Search(api, "\"je mange du\"")
    c.Search(api, "\"je mange des\"")
}

// Calls Twitter to execute the given query search.
// Stores the retrieved tweets into MongoDB, deduplicating
// them using their ID. Retweets aren't stored.
// Returns how many tweets were actually stored.
func (c *Crawler) Search(api *anaconda.TwitterApi, query string) int {
    // DAO
    tweetDao := db.NewTweetDAO(c.App.Mongo)

    // Maximum value for the search
    params := url.Values{}
    params.Set("count", "100")          // Amount of tweet possible in one query
    params.Set("result_type", "recent") // We want the more recent tweets
    searchResult, err := api.GetSearch(query, params)

    // Error, end of job for this time.
    if err != nil {
        log.Printf("An error ocurred during the search on Twitter : %s\n", err.Error())
        return 0
    }

    i := 0;

    // Saves every tweet into MongoDB for further analysis
    for _ , tweet := range searchResult {
        // Don't insert retweet
        if tweet.RetweetedStatus != nil {
            continue
        }

        // Don't insert it if we already have it (retweets)
        existing, err := tweetDao.FindByTweetId(tweet.Id)
        // Look for existing
        if len(existing) == 0 {
            err = tweetDao.Upsert(db.NewTweetFromApiTweet(&tweet))
            if err == nil {
                i++;
            }
        }
    }

    log.Printf("[info] %d tweets saved for \"%s\".\n", i, query);
    return i
}
