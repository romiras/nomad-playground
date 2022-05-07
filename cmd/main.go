package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/romiras/nomad-playground/internal/nomad"
)

var action = flag.String("action", "create", "A default action")

func main() {
	flag.Parse()
	nomad := nomad.NewNomad("")

	if action == nil {
		log.Fatal("action is not set")
	}

	switch *action {
	case "create":
		if err := RegisterJob(nomad); err != nil {
			panic(err)
		}
		fmt.Println("Done: job registered")
	case "delete":
		if err := Deregister(nomad); err != nil {
			panic(err)
		}
		fmt.Println("Done: job deregistered")
	}
}

func RegisterJob(n *nomad.NomadService) error {
	nomadJob := nomad.NomadJob{
		ID:         "job-id",
		Name:       "job-name",
		Region:     "",
		Priority:   50,
		Datacenter: "dc1",
		TaskGroups: []*nomad.NomadTaskGroup{
			{
				Name: "task-group-1",
				Tasks: []nomad.NomadTask{
					{
						Name:   "redis6-A",
						Driver: "docker",
						Config: map[string]interface{}{
							"image": "redis:6-alpine",
							"ports": []string{"p-redis"},
						},
						EnvVars:   nil,
						Resources: nil,
					},
					{
						Name:   "random-logger",
						Driver: "docker",
						Config: map[string]interface{}{
							"image": "chentex/random-logger:latest",
							"args": []string{
								"100", "400",
							},
						},
						EnvVars:   nil,
						Resources: nil,
					},
				},
			},
		},
	}

	job, err := n.Prepare(&nomadJob)
	if err != nil {
		return err
	}

	err = n.Register(job)
	if err != nil {
		return err
	}

	return nil
}

func Deregister(nomad *nomad.NomadService) error {
	return nomad.Deregister("job-id", false)
}

func sampleTaskGroup() *nomad.NomadTaskGroup {
	return &nomad.NomadTaskGroup{
		Name: "sample-task-group",
		Tasks: []nomad.NomadTask{
			{
				Name:   "Alpine",
				Driver: "docker",
				Config: map[string]interface{}{
					// "hostname":       "",
					"image": "alpine",
					// "command":        "",
					// "entrypoint":     []string{},
					// "args":           []string{}, // https://www.nomadproject.io/docs/runtime/interpolation
					// "volumes":        []string{},
					// "labels":         map[string]string{},
					// "network_mode":   "",
				},
				EnvVars:   map[string]string{},
				Resources: nil,
			},
		},
	}
}
