package task

import (
	"errors"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// tasks are one of the basic computing unite of an orchestrator
type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPort   nat.PortSet
	PortBinding   map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
	ContainerID   string
}

func (t *Task) GetState(state string) (State, error) {
	switch state {
	case Completed.Slug:
		return Completed, nil
	case Scheduled.Slug:
		return Scheduled, nil
	case Running.Slug:
		return Running, nil
	case Failed.Slug:
		return Failed, nil
	case Pending.Slug:
		return Pending, nil
	}
	return Unknown, errors.New("Unknown Type: " + state)
}
