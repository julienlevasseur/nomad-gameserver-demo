package main

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
)

type Endpoints struct {
	PostJobEndpoint         endpoint.Endpoint
	GetJobsEndpoint         endpoint.Endpoint
	GetJobEndpoint          endpoint.Endpoint
	PostJobDispatchEndpoint endpoint.Endpoint
	DeleteJobEndpoint       endpoint.Endpoint
	GetServicesEndpoint     endpoint.Endpoint
	GetServiceEndpoint      endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostJobEndpoint:         MakePostJobEndpoint(s),
		GetJobsEndpoint:         MakeGetJobsEndpoint(s),
		GetJobEndpoint:          MakeGetJobEndpoint(s),
		PostJobDispatchEndpoint: MakePostJobDispatchEndpoint(s),
		DeleteJobEndpoint:       MakeDeleteJobEndpoint(s),
		GetServicesEndpoint:     MakeGetServicesEndpoint(s),
		GetServiceEndpoint:      MakeGetServiceEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a profilesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		PostJobEndpoint:         httptransport.NewClient("POST", tgt, encodePostJobRequest, decodePostJobResponse, options...).Endpoint(),
		GetJobEndpoint:          httptransport.NewClient("GET", tgt, encodeGetJobRequest, decodeGetJobResponse, options...).Endpoint(),
		PostJobDispatchEndpoint: httptransport.NewClient("POST", tgt, encodePostJobDispatchRequest, decodePostJobDispatchResponse, options...).Endpoint(),
		DeleteJobEndpoint:       httptransport.NewClient("DELETE", tgt, encodeDeleteJobRequest, decodeDeleteJobResponse, options...).Endpoint(),
		GetServicesEndpoint:     httptransport.NewClient("GET", tgt, encodeGetServicesRequest, decodeGetServicesResponse, options...).Endpoint(),
		GetServiceEndpoint:      httptransport.NewClient("GET", tgt, encodeGetServiceRequest, decodeGetServiceResponse, options...).Endpoint(),
	}, nil
}

// PostJob implements Service. Primarily useful in a client.
func (e Endpoints) PostJob(ctx context.Context, j *nomad.Job) (*nomad.JobRegisterResponse, error) {
	request := postJobRequest{Job: j}
	response, err := e.PostJobEndpoint(ctx, request)
	if err != nil {
		return &nomad.JobRegisterResponse{}, err
	}
	resp := response.(postJobResponse)
	return resp.JobRegisterResponse, resp.Err
}

// GetJobs implements Service.
func (e Endpoints) GetJobs(ctx context.Context) ([]*nomad.JobListStub, error) {
	request := getJobsRequest{}
	response, err := e.GetJobsEndpoint(ctx, request)
	if err != nil {
		return []*nomad.JobListStub{}, err
	}
	resp := response.(getJobsResponse)
	return resp.Jobs, resp.Err
}

// GetJob implements Service.
func (e Endpoints) GetJob(ctx context.Context) (*nomad.Job, error) {
	request := getJobRequest{}
	response, err := e.GetJobEndpoint(ctx, request)
	if err != nil {
		return &nomad.Job{}, err
	}
	resp := response.(getJobResponse)
	return resp.Job, resp.Err
}

// PostJobDispatch implements Service. Primarily useful in a client.
func (e Endpoints) PostJobDispatch(ctx context.Context, id string, j DispatchJob) (*nomad.JobDispatchResponse, error) {
	request := postJobDispatchRequest{ID: id, Job: j}
	response, err := e.PostJobDispatchEndpoint(ctx, request)
	if err != nil {
		return &nomad.JobDispatchResponse{}, err
	}
	resp := response.(postJobDispatchResponse)
	return resp.JobDispatchResponse, resp.Err
}

// DeleteJob implements Service. Primarily useful in a client.
func (e Endpoints) DeleteJob(ctx context.Context, j DispatchJob) error {
	request := deleteJobRequest{Job: j}
	response, err := e.DeleteJobEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteJobResponse)
	return resp.Err
}

// MakeGetJobsEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetJobsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getJobsRequest)
		jobs, e := s.GetJobs(ctx)
		return getJobsResponse{Jobs: jobs, Err: e}, nil
	}
}

// MakeGetJobEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getJobRequest)
		job, e := s.GetJob(ctx, req.ID)
		return getJobResponse{Job: job, Err: e}, nil
	}
}

// MakePostJobEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postJobRequest)
		JobRegisterResponse, e := s.PostJob(ctx, req.Job)
		return postJobResponse{JobRegisterResponse: JobRegisterResponse, Err: e}, nil
	}
}

// MakePostJobDispatchEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostJobDispatchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postJobDispatchRequest)
		jdr, e := s.PostJobDispatch(ctx, req.ID, req.Job)
		return postJobDispatchResponse{JobDispatchResponse: jdr, Err: e}, nil
	}
}

// MakeDeleteJobEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteJobEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteJobRequest)
		e := s.DeleteJob(ctx, req.Job)
		return deleteJobResponse{Err: e}, nil
	}
}

// MakeGetServiceEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetServiceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getServiceRequest)
		services, e := s.GetService(ctx, req.ID)
		return getServiceResponse{Services: services, Err: e}, nil
	}
}

// MakeGetServicesEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetServicesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(getServicesRequest)
		services, e := s.GetServices(ctx)
		return getServicesResponse{Services: services, Err: e}, nil
	}
}

type postJobRequest struct {
	Job *nomad.Job
}

type postJobResponse struct {
	Err                 error                      `json:"err,omitempty"`
	JobRegisterResponse *nomad.JobRegisterResponse `json:"JobRegisterResponse"`
}

type getJobsRequest struct{}

type getJobsResponse struct {
	Err  error                `json:"err,omitempty"`
	Jobs []*nomad.JobListStub `json:"Jobs"`
}

type getJobRequest struct {
	ID string
}

type getJobResponse struct {
	Err error      `json:"err,omitempty"`
	Job *nomad.Job `json:"Job"`
}

func (r postJobResponse) error() error { return r.Err }

type postJobDispatchRequest struct {
	ID  string
	Job DispatchJob
}

type postJobDispatchResponse struct {
	Err                 error                      `json:"err,omitempty"`
	JobDispatchResponse *nomad.JobDispatchResponse `json:"JobDispatchResponse"`
}

type deleteJobRequest struct {
	Job DispatchJob
}

type deleteJobResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteJobResponse) error() error { return r.Err }

type getServiceRequest struct {
	ID string
}

type getServiceResponse struct {
	Err      error                    `json:"err,omitempty"`
	Services []*consul.CatalogService `json:"Services"`
}

type getServicesRequest struct{}

type getServicesResponse struct {
	Err      error               `json:"err,omitempty"`
	Services map[string][]string `json:"Services"`
}
