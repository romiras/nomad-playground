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

func RegisterJob(nomad *nomad.NomadService) error {
	job, err := nomad.Prepare("job-id", "job-name", 0)
	if err != nil {
		return err
	}

	err = nomad.Register(job)
	if err != nil {
		return err
	}

	return nil
}

func Deregister(nomad *nomad.NomadService) error {
	return nomad.Deregister("job-id", false)
}
