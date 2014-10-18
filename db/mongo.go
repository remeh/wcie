package db

import (
    "gopkg.in/v2/mgo"
)

// ---------------------- 
// Declarations

// A Mongo connection
// @author Rémy MATHIEU
type Mongo struct {
    session *mgo.Session
    database *mgo.Database
}

// ---------------------- 
// Methods

// Retrieves a new Mongo Connection.
// TODO pool of a sessions / connections
func GetConnection(mongoURI string) *Mongo {
    m := new(Mongo)
    session, err := mgo.Dial(mongoURI)
    if err != nil {
        panic(err)
    }
    m.session = session
    return m
}

func (m *Mongo) GetCollection(name string) *mgo.Collection {
    // TODO db name in configuration
    return m.session.DB("wcie").C(name)
}

func (m *Mongo) Close() {
    m.session.Close()
}
