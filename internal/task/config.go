package task

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// config needed to run a task
type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	Cmd           []string
	Image         string
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
}

// Docker client needed to interact with Docker.
type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerId string
}

// Return value to start and stop docker
type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

// Runs the method on the Docker client object, passing the context object,Run ImagePull
// the image name, and any options necessary to pull the image.
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: "hergytchuinkou",
		Auth:     "dckr_pat_P9bwidnItqV2zXcS0BzAOsRFSTA",
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	reader, err := d.Client.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)
	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory: d.Config.Memory,
	}

	cc := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	// The method earlier, returns two values, a response, whichImagePull ContainerCreate
	// is a pointer to a type, and an error type. Thecontainer.ContainerCreateCreatedBod.
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	if err := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting container using image %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	// printing stated container infos
	//d.Config.Runtime.ContainerID = resp.ID
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{ContainerId: resp.ID, Action: "start", Result: "success"}
}

// Docker client's container create methode based on given configuration
// func (cli *Docker) ContainerCreate(
// 	ctx context.Context,
// 	config *container.Config,
// 	hostConfig *container.HostConfig,
// 	networkingConfig *network.NetworkingConfig,
// 	platform *platform.Platform,
// 	containerName string) (container.ContainerCreateCreatedBody, error)

// The container stopped by calling the containerStop method,
// and finally it's removed by calling ContainerRemove.
func (d *Docker) Stop() DockerResult {
	ctx := context.Background()
	timeOut := 0
	log.Printf("Attempting to stop container %v", d.ContainerId)
	if err := d.Client.ContainerStop(ctx, d.ContainerId, container.StopOptions{Signal: "SIGTERM", Timeout: &timeOut}); err != nil {
		fmt.Printf("stopping container failed %s\n", err)

		//		panic(err)
	}
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}
	err := d.Client.ContainerRemove(ctx, d.ContainerId, removeOptions)
	if err != nil {
		fmt.Printf("stopping container failed %s\n", err)
		//		panic(err)
	}
	return DockerResult{Action: "stop", Result: "success", Error: nil}
}

func (t *Task) NewConfig(task *Task) *Config {
	return &Config{
		Name:         task.Name,
		Image:        task.Image,
		Memory:       134217728,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Disk:         64,
	}
}

func (t *Task) NewDocker(c *Config, client *client.Client) Docker {
	return Docker{
		Config:      *c,
		ContainerId: t.ContainerID,
		Client:      client,
	}
}
