package runtime

import (
    "fmt"
    "net/url"

    "github.com/ChimeraCoder/anaconda"
)

func Crawl(app *App) {
    // Api twitter provided by ChimeraCoder !
    anaconda.SetConsumerKey(app.Config.TwitterApiKey)
    anaconda.SetConsumerSecret(app.Config.TwitterSecret)
    api := anaconda.NewTwitterApi(app.Config.TwitterAccessToken, app.Config.TwitterAccessTokenSecret)
    defer api.Close()

    // Maximum value for the search
    params := url.Values{}
    params.Set("count", "100")
    searchResult, _ := api.GetSearch("\"je mange un\"", params)

    for _ , tweet := range searchResult {
        fmt.Println(tweet.Id)
        fmt.Println(tweet.CreatedAtTime())
        fmt.Println(tweet.Text)
    }
}
