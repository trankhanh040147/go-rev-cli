package ui

import "time"

// UI feedback durations
const (
	YankFeedbackDuration      = 2 * time.Second
	PruneErrorFeedbackDuration = 3 * time.Second
	YankChordTimeout          = 300 * time.Millisecond
)

