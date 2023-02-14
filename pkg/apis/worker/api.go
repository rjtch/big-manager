package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rjtch/big-manager/internal/task"
	"github.com/rjtch/big-manager/internal/worker"
)

type Api struct {
	Address string
	Port    int
	Worker  *worker.Worker
	Router  *chi.Mux
}

func (api *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	payload := json.NewDecoder(r.Body)
	payload.DisallowUnknownFields()
	event := task.TaskEvent{}
	err := payload.Decode(&event)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling paylod : %v\n", err)
		log.Printf("payload object : %v\n", event)
		log.Printf("Message : %v\n", msg)
		w.WriteHeader(400)
		// send response back to the CLI
		errRsp := ErrorResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(errRsp)
		return
	}
	// added task to the worker-queue
	api.Worker.AddTask(event.Task)
	log.Printf("Added task %v\n", event.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(event.Task)
}

func (api *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(api.Worker.GetTask())
}

func (api *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Printf("No taskID passed in request.\n")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	taskToStop := api.Worker.Db[tID]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	api.Worker.AddTask(taskCopy)
	_, ok := api.Worker.Db[tID]
	if ok {
		log.Printf("No task with ID %v found", tID)
		w.WriteHeader(404)
	}

	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(204)
}

func (api *Api) initRouter() {
	api.Router = chi.NewRouter()
	api.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", api.StartTaskHandler)
		r.Get("/", api.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", api.StopTaskHandler)
		})
	})
}

func (api *Api) Start() {
	api.initRouter()
	http.ListenAndServe("localhost:3333", api.Router)
}
