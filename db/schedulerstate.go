package db

import (
    "time"
)

// ----------------------
// Declarations

// Used to save the scheduler state.
type SchedulerState struct {
    Id string `bson:"_id"`
    LastScheduledMinute time.Time `bson:"last_scheduled_minute"`
}

// ----------------------
// Methods

// Constructor
func NewSchedulerState() *SchedulerState {
    var t time.Time
    return &SchedulerState{
        Id: "state_v1",
        LastScheduledMinute: t,
    }
}

