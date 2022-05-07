package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type NomadService struct {
	client *api.Client
	region string
}

const (
	DefaultPriority = 50
)

func NewNomad(region string) *NomadService {
	config := api.DefaultConfig()

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	if region == "" {
		region = api.GlobalRegion
	}

	fmt.Printf("client: %#v\n", client)

	return &NomadService{
		client: client,
		region: region,
	}
}

func (n *NomadService) Prepare(jobID, jobName string, priority int) (*api.Job, error) {
	if priority == 0 {
		priority = DefaultPriority
	}
	job := api.NewServiceJob(jobID, jobName, n.region, priority)
	job.AddDatacenter("dc1")

	tasks := []*api.Task{
		n.createTask("random-logger", "docker", nil, map[string]interface{}{
			"image": "chentex/random-logger:latest",
			"args": []string{
				"100", "400",
			},
		}),
	}
	taskGroup := n.createTaskGroup("task-group-name", tasks)

	job.AddTaskGroup(taskGroup)

	if len(job.TaskGroups) == 0 {
		panic("no task groups")
	}

	var err error
	var writeMeta *api.WriteMeta
	var valResp *api.JobValidateResponse
	var planResp *api.JobPlanResponse

	jobs := n.client.Jobs()
	valResp, writeMeta, err = jobs.Validate(job, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("valResp: %#v, writeMeta: %#v\n", valResp, writeMeta)

	planResp, writeMeta, err = jobs.Plan(job, true, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("planResp: %#v, writeMeta: %#v\n", planResp, writeMeta)

	return job, nil
}

func (n *NomadService) Plan(job *api.Job) error {
	jobs := n.client.Jobs()
	planResponse, writeMeta, err := jobs.Plan(job, true, nil)
	if err != nil {
		return err
	}
	fmt.Printf("planResponse: %#v, writeMeta: %#v\n", planResponse, writeMeta)

	return nil
}

func (n *NomadService) Register(job *api.Job) error {
	regResponse, writeMeta, err := n.client.Jobs().Register(job, nil)
	if err != nil {
		return err
	}
	fmt.Printf("regResponse: %#v, writeMeta: %#v\n", regResponse, writeMeta)

	return nil
}

func (n *NomadService) Deregister(jobID string, purge bool) error {
	deregResponse, writeMeta, err := n.client.Jobs().Deregister(jobID, purge, nil)
	if err != nil {
		return err
	}
	fmt.Printf("deregResponse: %#v, writeMeta: %#v\n", deregResponse, writeMeta)

	return nil
}

func (n *NomadService) createTaskGroup(taskGroupName string, tasks []*api.Task) *api.TaskGroup {
	taskGroup := api.NewTaskGroup(taskGroupName, 1)

	for _, task := range tasks {
		taskGroup.AddTask(task)
	}

	return taskGroup
}

func (n *NomadService) createTask(taskName, taskDriver string, envVars map[string]string, config map[string]interface{}) *api.Task {
	task := api.NewTask(taskName, taskDriver)
	task.Env = envVars
	task.Config = config
	task.Require(&api.Resources{
		CPU:      intToPtr(100),
		MemoryMB: intToPtr(256),
	})
	return task
}
