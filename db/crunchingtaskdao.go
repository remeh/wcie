package db

import (
	"gopkg.in/v2/mgo"
	"gopkg.in/v2/mgo/bson"
)

// ----------------------
// Declarations

// DAO for crunching task collection.
// @author Rémy MATHIEU
type CrunchingTaskDAO struct {
	mongo      *Mongo
	collection *mgo.Collection
}

// ----------------------
// Methods

func NewCrunchingTaskDAO(m *Mongo) *CrunchingTaskDAO {
	return &CrunchingTaskDAO{m, m.GetCollection(C_CRUNCHING_TASK)}
}

// Save/update the given tweet.
func (d *CrunchingTaskDAO) Upsert(task *CrunchingTask) error {
    _, err := d.collection.Upsert(bson.M{"_id": task.Id}, task)
    return err
}

// Gets a task to compute, removing it from the list.
func (d *CrunchingTaskDAO) GetNext(limit int) ([]CrunchingTask, error) {
	var q *mgo.Query
	q = d.collection.Find(bson.M{}).Limit(limit)

	iter := q.Iter()
	var task CrunchingTask
	tasks := make([]CrunchingTask, 0)
	for iter.Next(&task) {
		tasks = append(tasks, task)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

    return tasks, nil
}

// Deletes the given task from MongoDB
func (d *CrunchingTaskDAO) RemoveAll(tasks []CrunchingTask) error {
    // FIXME Bulk delete
    for _, task := range tasks {
        err := d.collection.Remove(task)
        if err != nil {
            return err
        }
    }
    return nil
}
