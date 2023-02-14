package task

import (
	"time"

	"github.com/google/uuid"
)

// eventing is used to inform a manager about any new instruction for a task
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}
