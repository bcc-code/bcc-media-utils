package jobs

import "github.com/google/uuid"

type Job struct {
	ID          uuid.UUID
	Query       string
	Status      string
	ItemIDs     []string
	Concurrency int
}

type JobItem struct {
	ID     uuid.UUID
	ItemID string
	Status string
}
