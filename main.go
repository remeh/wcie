// Launches the crawl job,
// the HTTP webserver and
// the computing system.
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/remeh/wcie/runtime"
)

func main() {
    app := runtime.NewApp(prepareFlags())
    app.Start()
}

// Prepares the CLI flags for the
// Mongo connection and the file to import.
func prepareFlags() runtime.Config {
    mongoURI := flag.String("m", "localhost", "The Mongo URI to connect to MongoDB.")
    dbName := flag.String("d", "wcie", "The DB name to use in MongoDB.")
    apiKey := flag.String("k", "", "The twitter apikey to use for queries on Twitter.")
    secret := flag.String("s", "", "The twitter secret use as a password on Twitter to identify this app.")
    token := flag.String("t", "", "The access token to query Twitter.")
    tokenSecret := flag.String("ts", "", "The access token secret to query Twitter.")

    flag.Parse()

    // Mandatory parameters.
    if len(*apiKey) == 0 || len(*secret) == 0 || len(*token) == 0 || len(*tokenSecret) == 0 {
        fmt.Println("Bad usage: one or more missing flags. See :")
        flag.PrintDefaults()
        os.Exit(1)
    }

    return runtime.Config{MongoURI: *mongoURI, DBName: *dbName, TwitterApiKey: *apiKey, TwitterSecret: *secret, TwitterAccessToken: *token, TwitterAccessTokenSecret: *tokenSecret}
}

