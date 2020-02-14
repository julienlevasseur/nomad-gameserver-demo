package main

import (
	"context"

	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
)

// Service is a simple CRUD interface for Nomad Jobs & Consul Services.
type Service interface {
	PostJob(ctx context.Context, j *nomad.Job) (*nomad.JobRegisterResponse, error)
	GetJobs(ctx context.Context) ([]*nomad.JobListStub, error)
	GetJob(ctx context.Context, id string) (*nomad.Job, error)
	PostJobDispatch(ctx context.Context, id string, j DispatchJob) (*nomad.JobDispatchResponse, error)
	DeleteJob(ctx context.Context, j DispatchJob) error
	GetServices(ctx context.Context) (map[string][]string, error)
	GetService(ctx context.Context, id string) ([]*consul.CatalogService, error)
}

// DispatchJob is a job with parameters
type DispatchJob struct {
	JobID   string            `json:"JobID"`
	Meta    map[string]string `json:"Meta"`
	Payload []byte            `json:"Payload"`
}

type api struct{}

// NewAPI return a new instance of API
func newAPI() *api {
	return &api{}
}

func nomadClient() (*nomad.Jobs, error) {
	n, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		return &nomad.Jobs{}, err
	}
	return n.Jobs(), nil
}

func consulCatalog() (*consul.Catalog, error) {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return consulClient.Catalog(), nil
}

// PostJob implements the Nomad Job Register call.
func (a *api) PostJob(ctx context.Context, job *nomad.Job) (*nomad.JobRegisterResponse, error) {
	n, err := nomadClient()
	if err != nil {
		return &nomad.JobRegisterResponse{}, err
	}

	jobRegisterResponse, _, err := n.Register(job, nil)
	if err != nil {
		return &nomad.JobRegisterResponse{}, err
	}
	return jobRegisterResponse, nil
}

// GetJobs implements the Nomad Job List call.
func (a *api) GetJobs(ctx context.Context) ([]*nomad.JobListStub, error) {
	n, err := nomadClient()
	if err != nil {
		return []*nomad.JobListStub{}, err
	}

	jobs, _, err := n.List(nil)
	if err != nil {
		return []*nomad.JobListStub{}, err
	}

	return jobs, nil
}

// GetJob implements the Nomad Job Info call.
func (a *api) GetJob(ctx context.Context, id string) (*nomad.Job, error) {
	n, err := nomadClient()
	if err != nil {
		return &nomad.Job{}, err
	}

	job, _, err := n.Info(id, nil)
	if err != nil {
		return &nomad.Job{}, err
	}

	return job, nil
}

// PostJobDispatch implements the Nomad Job Dispatch call.
func (a *api) PostJobDispatch(ctx context.Context, id string, j DispatchJob) (*nomad.JobDispatchResponse, error) {
	n, err := nomadClient()
	if err != nil {
		return &nomad.JobDispatchResponse{}, err
	}

	jobDispatchResponse, _, err := n.Dispatch(id, j.Meta, j.Payload, nil)
	if err != nil {
		return &nomad.JobDispatchResponse{}, err
	}

	return jobDispatchResponse, nil
}

// DeleteJob implements the Nomad Job Deregister call.
func (a *api) DeleteJob(ctx context.Context, j DispatchJob) error {
	n, err := nomadClient()
	if err != nil {
		return err
	}

	_, _, err = n.Deregister(j.JobID, false, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetServices implements the Consul Catalog Services call.
func (a *api) GetServices(ctx context.Context) (map[string][]string, error) {
	catalog, err := consulCatalog()
	if err != nil {
		return make(map[string][]string), err
	}

	services, _, err := catalog.Services(nil)
	if err != nil {
		return make(map[string][]string), err
	}

	return services, nil
}

// GetService implements the Consul Catalog Service call.
func (a *api) GetService(ctx context.Context, id string) ([]*consul.CatalogService, error) {
	catalog, err := consulCatalog()
	if err != nil {
		return []*consul.CatalogService{}, err
	}

	catalogService, _, err := catalog.Service(id, "", nil)
	if err != nil {
		return []*consul.CatalogService{}, err
	}

	return catalogService, nil
}
