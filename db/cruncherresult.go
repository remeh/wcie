package db

import (
    "time"
)

// Comptued results
type CruncherResult struct {
    Id time.Time `bson:"_id"` // The ID representing the time computed. Second will be to 1 if it's the hour computing.

    Data map[string]int
}
