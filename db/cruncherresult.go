package db

import (
    "time"
)

// Computed results
type CruncherResult struct {
    Id time.Time `bson:"_id"` // The ID representing the time computed. Second will be to 1 if it's the hour computing.

    CreationTime time.Time `bson:"creation_time"` // At which time this result has been computed

    Data map[string]int
}
