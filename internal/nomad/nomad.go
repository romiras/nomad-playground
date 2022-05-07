package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type NomadService struct {
	client *api.Client
	region string
}

type (
	NomadJob struct {
		ID         string
		Name       string
		Datacenter string
		Region     string
		Priority   int
		TaskGroups []*NomadTaskGroup
	}

	NomadTaskGroup struct {
		Name  string
		Tasks []NomadTask
	}

	NomadTask struct {
		Name      string
		Driver    string
		Config    map[string]interface{}
		EnvVars   map[string]string
		Resources *NomadTaskResources
	}

	NomadTaskResources struct {
		CPU         *int
		Cores       *int
		MemoryMB    *int
		MemoryMaxMB *int
		DiskMB      *int
	}
)

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

func (n *NomadService) Prepare(nomadJob *NomadJob) (*api.Job, error) {
	if nomadJob == nil {
		return nil, fmt.Errorf("nomadJob is nil")
	}

	if nomadJob.Priority == 0 {
		nomadJob.Priority = DefaultPriority
	}
	job := api.NewServiceJob(nomadJob.ID, nomadJob.Name, nomadJob.Region, nomadJob.Priority)
	job.AddDatacenter(nomadJob.Datacenter)

	for _, tg := range nomadJob.TaskGroups {
		tasks := make([]*api.Task, 0, len(tg.Tasks))
		for _, nomadTask := range tg.Tasks {
			tasks = append(tasks, n.createTask(&nomadTask))
		}
		taskGroup := n.createTaskGroup(tg.Name, tasks)
		taskGroup.Canonicalize(job)
		job.AddTaskGroup(taskGroup)
	}

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
	taskGroup.Networks = []*api.NetworkResource{
		{
			DynamicPorts: []api.Port{
				{Label: "p-redis", Value: 0, To: 6379},
			},
		},
	}

	for _, task := range tasks {
		taskGroup.AddTask(task)
	}

	return taskGroup
}

func (n *NomadService) createTask(nomadTask *NomadTask) *api.Task {
	task := api.NewTask(nomadTask.Name, nomadTask.Driver)
	task.Env = nomadTask.EnvVars
	task.Config = nomadTask.Config
	if nomadTask.Resources != nil {
		res := api.Resources{}
		if nomadTask.Resources.CPU != nil {
			res.CPU = intToPtr(*nomadTask.Resources.CPU)
		}
		if nomadTask.Resources.Cores != nil {
			res.Cores = intToPtr(*nomadTask.Resources.Cores)
		}
		if nomadTask.Resources.MemoryMB != nil {
			res.MemoryMB = intToPtr(*nomadTask.Resources.MemoryMB)
		}
		if nomadTask.Resources.MemoryMaxMB != nil {
			res.MemoryMaxMB = intToPtr(*nomadTask.Resources.MemoryMaxMB)
		}
		if nomadTask.Resources.DiskMB != nil {
			res.DiskMB = intToPtr(*nomadTask.Resources.DiskMB)
		}
		task.Require(&res)
	} else {
		task.Require(&api.Resources{
			CPU:      intToPtr(100),
			MemoryMB: intToPtr(256),
		})
	}

	return task
}
