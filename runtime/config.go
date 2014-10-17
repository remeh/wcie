package runtime

// App configuration.
type Config struct {
    MongoURI string // the mongo connection string
    DBName string // db to use in MongoDB
    TwitterApiKey string // The twitter apikey
    TwitterSecret string // Secret to use to identify on twitter this app
    TwitterAccessToken string // The access token to query Twitter
    TwitterAccessTokenSecret string // The access token secret to query Twitter
}
