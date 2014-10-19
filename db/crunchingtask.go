package db

import (
    "time"
)

// A crunching task is eaten by the cruncher
// to know what it has to compute.
type CrunchingTask struct {
    Id int64 `bson:"_id"` // The ID is in the format : 201412200100 to eliminate duplicate
    CreationTime time.Time `bson:"creation_time"` // At which time this task has been created.
}
