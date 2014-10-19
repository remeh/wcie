package db

import (
	"gopkg.in/v2/mgo"
	"gopkg.in/v2/mgo/bson"
)

// ----------------------
// Declarations

type SchedulerStateDAO struct {
	mongo      *Mongo
	collection *mgo.Collection
}

// ----------------------
// Methods

func NewSchedulerStateDAO(m *Mongo) *SchedulerStateDAO {
	return &SchedulerStateDAO{m, m.GetCollection(C_SCHEDULER_STATE)}
}

// Retrieves the scheduler state from the DB.
func (d *SchedulerStateDAO) Save(schedulerState *SchedulerState) error {
    _, err := d.collection.Upsert(bson.M{"_id": schedulerState.Id}, schedulerState)
    return err
}

// Retrieves the scheduler state from the DB.
func (d *SchedulerStateDAO) Get() *SchedulerState {
	var q *mgo.Query

	q = d.mongo.GetCollection(C_SCHEDULER_STATE).Find(bson.M{})

	iter := q.Iter()
	var state SchedulerState
	states := make([]*SchedulerState, 0)
	for iter.Next(&state) {
		states = append(states, &state)
	}

    count := len(states)
    switch {
        case count > 1:
            // Many scheduler states found, we must stop the application
            // right now and some states should be retrieved manually
            // from the database.
            panic("Too many scheduler states found in the database.")
        case count == 1:
            return states[0]
    }

    // No scheduler found, create one.
    return NewSchedulerState()
}
