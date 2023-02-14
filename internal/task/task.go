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
	case Completed.slug:
		return Completed, nil
	case Scheduled.slug:
		return Scheduled, nil
	case Running.slug:
		return Running, nil
	case Failed.slug:
		return Failed, nil
	case Pending.slug:
		return Pending, nil
	}
	return Unknown, errors.New("Unknown Type: " + state)
}
