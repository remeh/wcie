package db

import (
	"gopkg.in/v2/mgo"
	"gopkg.in/v2/mgo/bson"
)

// ----------------------
// Declarations

// DAO for computed cruncher results.
// @author Rémy MATHIEU
type CruncherResultDAO struct {
	mongo      *Mongo
	collection *mgo.Collection
}

// ----------------------
// Methods

func NewCruncherResultDAO(m *Mongo) *CruncherResultDAO {
	return &CruncherResultDAO{m, m.GetCollection(C_CRUNCHER_RESULT)}
}

// Save/update the given results.
func (d *CruncherResultDAO) Upsert(task *CruncherResult) error {
    _, err := d.collection.Upsert(bson.M{"_id": task.Id}, task)
    return err
}

