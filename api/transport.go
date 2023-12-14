package main

// The profilesvc is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func (s *Service) MakeHTTPHandler(logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := s.MakeServerEndpoints()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/v1/jobs/").Handler(httptransport.NewServer(
		e.PostJobEndpoint,
		decodePostJobRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/v1/jobs/").Handler(httptransport.NewServer(
		e.GetJobsEndpoint,
		decodeGetJobsRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/v1/jobs/").Handler(httptransport.NewServer(
		e.DeleteJobEndpoint,
		decodeDeleteJobRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/v1/jobs/{id}").Handler(httptransport.NewServer(
		e.GetJobEndpoint,
		decodeGetJobRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/v1/jobs/{id}/dispatch").Handler(httptransport.NewServer(
		e.PostJobDispatchEndpoint,
		decodePostJobDispatchRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/v1/services/").Handler(httptransport.NewServer(
		e.GetServicesEndpoint,
		decodeGetServicesRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/v1/services/{id}").Handler(httptransport.NewServer(
		e.GetServiceEndpoint,
		decodeGetServiceRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostJobRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var job *nomad.Job
	job, err = jobspec.Parse(r.Body)
	return postJobRequest{
		Job: job,
	}, nil
}

func decodeGetJobsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getJobsRequest{}, nil
}

func decodeGetJobRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getJobRequest{
		ID: id,
	}, nil
}

func decodePostJobDispatchRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var job DispatchJob
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		return nil, err
	}
	return postJobDispatchRequest{
		ID:  id,
		Job: job,
	}, nil
}

func decodeDeleteJobRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var job DispatchJob
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		return nil, err
	}
	return deleteJobRequest{
		Job: job,
	}, nil
}

func decodeGetServicesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getServicesRequest{}, nil
}

func decodeGetServiceRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getServiceRequest{
		ID: id,
	}, nil
}

func encodePostJobRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/jobs/"
	return encodeRequest(ctx, req, request)
}

func encodeGetJobsRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/jobs/"
	return encodeRequest(ctx, req, request)
}

func encodeGetJobRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getJobRequest)
	jobID := url.QueryEscape(r.ID)
	req.URL.Path = "/jobs/" + jobID
	return encodeRequest(ctx, req, request)
}

func encodePostJobDispatchRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(postJobDispatchRequest)
	jobID := url.QueryEscape(r.ID)
	req.URL.Path = "/jobs/" + jobID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteJobRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/jobs/"
	return encodeRequest(ctx, req, request)
}

func encodeGetServicesRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/services/"
	return encodeRequest(ctx, req, request)
}

func encodeGetServiceRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getServiceRequest)
	serviceID := url.QueryEscape(r.ID)
	req.URL.Path = "/services/" + serviceID
	return encodeRequest(ctx, req, request)
}

func decodePostJobResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postJobResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetJobResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getJobResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostJobDispatchResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postJobDispatchResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteJobResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteJobResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetServicesResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getServicesResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetServiceResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getServiceResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
