package worker

import (
	"errors"
	"fmt"
	"github.com/rjtch/big-manager/pkg/apis/metrics/linux"
	"log"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/rjtch/big-manager/internal/task"
)

type Worker struct {
	Queue     *queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
	Name      string
	Stats     *linux.Stats
}

func (w *Worker) CollectStats() {
	for {
		fmt.Println("I will collect stats")
		w.Stats = linux.GetStats()
		time.Sleep(30)
	}
}

func (w *Worker) RunTask(client *client.Client) task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Printf("No tasks present in the queue")
		return task.DockerResult{Error: nil}
	}
	// convert task from the queue to our task object since dequeue returns
	// object from type interface
	taskQueue := t.(task.Task)
	// get the same task from the Db
	taskPersisted := w.Db[taskQueue.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueue
		w.Db[taskPersisted.ID] = &taskQueue
	}
	//validate state of the task
	var result task.DockerResult
	if ValidStateTransition(taskPersisted.State, taskQueue.State) {
		switch taskQueue.State {
		case task.Scheduled:
			return w.StartTask(taskQueue, client)
		case task.Completed:
			return w.StopTask(taskQueue, client)
		default:
			// the task has either failed or does not exist.
			result.Error = errors.New("The start you tried to start is does not exist")
		}
	} else {
		err := fmt.Errorf("Invalid transmition from %v to %v", taskPersisted.State, taskQueue.State)
		result.Error = err
	}
	return result
}

func (w *Worker) StartTask(task task.Task, client *client.Client) task.DockerResult {
	config := task.NewConfig(&task)
	d := task.NewDocker(config, client)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Error starting container %v: %v\n", d.ContainerId, result.Error)
		task.State, _ = task.GetState("Failed")
		w.Db[task.ID] = &task
		return result
	}

	d.ContainerId = result.ContainerId
	task.State, _ = task.GetState("Running")
	w.Db[task.ID] = &task

	return result
}

func (w *Worker) StopTask(task task.Task, client *client.Client) task.DockerResult {
	config := task.NewConfig(&task)
	d := task.NewDocker(config, client)
	result := d.Stop()
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", d.ContainerId, result.Error)
	}
	task.FinishTime = time.Now().UTC()
	task.State, _ = task.GetState("Completed")
	w.Db[task.ID] = &task
	log.Printf("Stopped and removed container  %v for task %v", d.ContainerId, task.ID)

	return result
}

func (w *Worker) AddTask(t task.Task) {
	if w.Queue.Len() <= 0 {

	}
	w.Queue.Enqueue(t)
}

func (w *Worker) GetTask() []task.Task {
	var allTasks []task.Task
	for w.Queue.Len() > 0 {
		allTasks = append(allTasks, w.Queue.Peek().(task.Task))
		log.Printf("All tasks %v\n ", allTasks)
	}
	return allTasks
}

func Contains(states []task.State, state task.State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidStateTransition(src task.State, dst task.State) bool {
	return Contains(task.StateTransitionMap[src], dst)
}
