// Launches the crawl job,
// the HTTP webserver and
// the computing system.
package main

import (
    "flag"

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
    flag.Parse()

    return runtime.Config{MongoURI: *mongoURI, DBName: *dbName}
}

