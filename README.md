# What Can I Eat

WCIE is a basic side-project with just a few hours used for its development.

Its purpose is to gather tweets using some search queries on the Twitter API, after having stored them into MongoDB, its storage model allows to do aggregation on tweets wrote during a given interval of time. This storage model is used to "discover" the next word to the search query in order to basically analyze what people are used to eat every day. 
The results will obviously be approximative, precision isn't the purpose of this project.

WCIE is licensed with the Apache License 2.0, don't hesitate to send ideas, issues or anything.

## Usage

```
  -d="wcie": The DB name to use in MongoDB.
  -k="": The twitter apikey to use for queries on Twitter.
  -m="localhost": The Mongo URI to connect to MongoDB.
  -s="": The twitter secret use as a password on Twitter to identify this app.
  -t="": The access token to query Twitter.
  -ts="": The access token secret to query Twitter.
```

## Roadmap

  * Compute percentage of each words per days, weeks, months, ...
  * HTTP route to display (with D3.js ?) some charts.

## Dependencies

```
  github.com/ChimeraCoder/anaconda
  gopkg.in/v2/mgo
```
