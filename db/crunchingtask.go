package db

import (
    "time"
)

// A crunching task is eaten by the cruncher
// to know what it has to compute.
type CrunchingTask struct {
    Id time.Time `bson:"_id"` // The ID representing the time to compute. Second will be to 1 if we need to compute the full hour.
    CreationTime time.Time `bson:"creation_time"` // At which time this task has been created.
}
