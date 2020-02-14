package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostJob(ctx context.Context, j *nomad.Job) (jrr *nomad.JobRegisterResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostJob", "id", j, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostJob(ctx, j)
}

func (mw loggingMiddleware) GetJobs(ctx context.Context) (jobs []*nomad.JobListStub, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetJobs", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetJobs(ctx)
}

func (mw loggingMiddleware) GetJob(ctx context.Context, id string) (job *nomad.Job, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetJob", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetJob(ctx, id)
}

func (mw loggingMiddleware) DeleteJob(ctx context.Context, j DispatchJob) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteJob", "id", j.JobID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteJob(ctx, j)
}

func (mw loggingMiddleware) PostJobDispatch(ctx context.Context, id string, j DispatchJob) (jdr *nomad.JobDispatchResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostJobDispatch", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostJobDispatch(ctx, id, j)
}

func (mw loggingMiddleware) GetServices(ctx context.Context) (services map[string][]string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetServices", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetServices(ctx)
}

func (mw loggingMiddleware) GetService(ctx context.Context, id string) (services []*consul.CatalogService, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetService", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetService(ctx, id)
}
