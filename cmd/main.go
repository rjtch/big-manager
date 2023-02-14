package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/rjtch/big-manager/internal/task"
	"github.com/rjtch/big-manager/internal/worker"
	workers "github.com/rjtch/big-manager/pkg/apis/worker"
)

func main() {
	// test container creation
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	host := os.Getenv("BIG_HOST")
	port, _ := strconv.Atoi(os.Getenv("BIG_PORT"))

	fmt.Println("Starting Big manager worker")
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := workers.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w, cli)
	api.Start()
}

func runTasks(w *worker.Worker, client *client.Client) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask(client)
			if result.Error != nil {
				fmt.Sprintf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently\n")
		}
		log.Printf("Sleeping for 10 Seconds.")
		time.Sleep(10 * time.Second)
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test_container_1",
		Image: "alpine",
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	// create container
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config \n", result.Result)
	return &d, &result
}

func stopContainer(d *task.Docker) *task.DockerResult {
	result := d.Stop()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf("Container %s is been stopping\n", result.ContainerId)
	return &result
}
