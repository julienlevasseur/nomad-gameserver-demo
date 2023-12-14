package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/julienlevasseur/nomad-gameserver-demo/api"
)

type Service interface {
	MakeHTTPHandler(logger log.Logger) http.Handler
	MakeServerEndpoints() api.Endpoints
}

type service struct{}

func main() {
	//	var (
	//		httpAddr = flag.String("http.addr", ":7070", "HTTP listen address")
	//	)
	//	flag.Parse()
	//
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var s Service
	//
	//	{
	//		s = newAPI()
	//		s = LoggingMiddleware(logger)(s)
	//	}
	//
	var h http.Handler
	//
	//	{
	h = s.MakeHTTPHandler(logger)
	//	}
	//
	// errs := make(chan error)
	//
	//	go func() {
	//		c := make(chan os.Signal)
	//		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//		errs <- fmt.Errorf("%s", <-c)
	//	}()
	//
	//	go func() {
	//		logger.Log("transport", "HTTP", "addr", *httpAddr)
	//		errs <- http.ListenAndServe(*httpAddr, h)
	//	}()
	//
	// logger.Log("exit", <-errs)
}
